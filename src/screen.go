package src

import (
	"github.com/gdamore/tcell"
	"sync"
)

const (
	DEFAULT_WIDTH  = 200
	DEFAULT_HEIGHT = 100

	WORLD_VIEW     = 0x0
	SCREEN_VIEW    = 0x1

	Z_DEPTH        = 0x4
)

type InputData struct {
	MouseDepth     int
	MousePressed   rune
	MouseHeld      rune
	PrevMouse      rune
	KeyPressed     rune
	KeyHeld        rune
	PrevKey        rune
	// used for special key presses e.g. ctrl + c
	CtrlKeyPressed tcell.Key
	CtrlKeyHeld    tcell.Key
	PrevCtrlKey    tcell.Key
	MousePos       Vec
}

type Screen struct {
	Screen        tcell.Screen
	Cam    		  Vec
	Width, Height int
	ZBuffer       []tcell.CellBuffer
	PrevZBuffer   []tcell.CellBuffer
}

var (
	ScreenMutex  = sync.Mutex{}
	InputBuffer = InputData{}
)

func (Screen *Screen) Init(){
	screen, err := tcell.NewScreen()
	if err != nil{
		CLogErr(err)
	}
	Screen.Screen = screen
	if err = screen.Init(); err != nil{
		CLogErr(err)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.Clear()

	Screen.Width = DEFAULT_WIDTH
	Screen.Height = DEFAULT_HEIGHT
	Screen.Resize()
	Screen.ZBuffer = make([]tcell.CellBuffer, Z_DEPTH)
	for i:=0; i<Z_DEPTH; i++ {
		Screen.ZBuffer[i] = tcell.CellBuffer{}
		Screen.ZBuffer[i].Resize(Screen.Width, Screen.Height)
	}
	InputBuffer = InputData{
		MousePressed:   0,
		MouseHeld:      0,
		PrevMouse:      -1,
		KeyPressed:     0,
		KeyHeld:        0,
		PrevKey:        -1,
		CtrlKeyPressed: 0,
		MousePos:       Vec{},
	}
}

func (Screen *Screen) Resize(){
	for _, buf := range Screen.ZBuffer {
		buf.Resize(Screen.Width, Screen.Height)
	}
	Screen.Screen.Sync()
}

func (Screen *Screen) WorldToScreen(v Vec) Vec{
	vNew := v.Add(Screen.Cam)
	return vNew.Add(V2i(Screen.Width/2, Screen.Height/2))
}

func (Screen *Screen) ScreenToWorld(v Vec) Vec{
	vNew := v.Sub(V2i(Screen.Width/2, Screen.Height/2))
	return vNew.Sub(Screen.Cam)
}

// TODO We may be able to just write the text directly
func (Screen *Screen) Text(text string, pos Vec, style tcell.Style, view uint8, depth int){
	for i, r := range text {
		if int(pos.X)+i >= Screen.Width {
			pos.X=0
			pos.Y++
		}
		Screen.Char(r, V2i(int(pos.X)+i, int(pos.Y)), style, view, depth)
	}
}

func (Screen *Screen) Char(r rune, pos Vec, style tcell.Style, view uint8, depth int){
	if view == WORLD_VIEW{
		pos = Screen.WorldToScreen(pos)
	}
	Screen.ZBuffer[depth].SetContent(int(pos.X), int(pos.Y), r, nil, style)
}

func (Screen *Screen) Rect(r rune, pos Vec, width, height int, style tcell.Style, fill bool, view uint8, depth int){
	for xTmp:=int(pos.X); xTmp<int(pos.Y)+width; xTmp++ {
		for yTmp:=int(pos.X); yTmp<int(pos.Y)+height; yTmp++ {
			if fill {
				Screen.Char(r, V2i(xTmp, yTmp), style, view, depth)
			}else{
				if (xTmp ==int(pos.X) || xTmp==int(pos.X)+width-1) || (yTmp==int(pos.Y) || yTmp==int(pos.Y)+height-1) {
					Screen.Char(r, V2i(xTmp, yTmp), style, view, depth)
				}
			}
		}
	}
}

func (Screen *Screen) Draw(){
	ScreenMutex.Lock()

	// reset what button was pressed
	InputBuffer.MousePressed = 0
	InputBuffer.KeyPressed=0
	InputBuffer.KeyHeld=0
	InputBuffer.CtrlKeyPressed=0
	InputBuffer.CtrlKeyHeld=0

	for y := 0; y < Screen.Height; y++ {
		for x := 0; x < Screen.Width; x++ {
			for i:=0; i<Z_DEPTH;i++ {
				rune, _, style, _ := Screen.ZBuffer[i].GetContent(x, y)
				// we have to draw the cell if the z-depth is 0 to clear the screen,
				// otherwise we get drawing artifacts. The reason we draw only if the rune
				// if not ' ' is because otherwise it would clear the char in the lower depth.
				if i == 0 || rune != ' ' {
					Screen.Screen.SetCell(x, y, style, rune)
				}
			}
		}
	}
	Screen.Screen.Show()
	// first copy the zbuffer to the previous zbuffer
	Screen.PrevZBuffer = make([]tcell.CellBuffer, len(Screen.ZBuffer))
	copy(Screen.PrevZBuffer, Screen.ZBuffer)
	for i:=0; i<Z_DEPTH; i++ {
		Screen.ZBuffer[i] = tcell.CellBuffer{}
		Screen.ZBuffer[i].Resize(Screen.Width, Screen.Height)
	}
	ScreenMutex.Unlock()
}

// TODO the screen seems to be able to only recognise clicks on depth 0, every other z-layer is just blank
func (Screen *Screen) CalculateMouseDepth() int{
	depth := -1
	// go through the screen buffers at the position of the mouse to see what is being drawn
	for i:=0;i<Z_DEPTH;i++{
		// why is this prevzbuffer essentially 0?
		rune, _, _, _ := Screen.PrevZBuffer[i].GetContent(int(InputBuffer.MousePos.X), int(InputBuffer.MousePos.Y))
		// if the layer contains a rune at the position we requested, then set the depth to that layer
		if rune != ' '{
			depth = i
		}
	}
	return depth
}

func (Screen *Screen) Poll() {
	for Running {
		ev := Screen.Screen.PollEvent()
		// when we receive an event, lock the screen
		ScreenMutex.Lock()
		switch ev := ev.(type){
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune{
				key := ev.Rune()
				InputBuffer.PrevKey = InputBuffer.KeyHeld
				InputBuffer.KeyHeld = key
				// if we are pressing a different key, do a key press
				if InputBuffer.KeyHeld != InputBuffer.PrevKey{
					InputBuffer.KeyPressed = key
				}
			}else {
				key := ev.Key()
				InputBuffer.PrevCtrlKey = InputBuffer.CtrlKeyHeld
				InputBuffer.CtrlKeyHeld = key
				// if we are pressing a different key, do a key press
				if InputBuffer.CtrlKeyHeld != InputBuffer.PrevCtrlKey{
					InputBuffer.CtrlKeyPressed = key
				}
			}
			break
		case *tcell.EventMouse:
			switch ev.Buttons(){
			// mouse released, we set the button pressed to what was being held down
			case tcell.ButtonNone:
				// if we are actually holding the mouse down, then its a valid click release
				if InputBuffer.MouseHeld != 0 {
					InputBuffer.MousePressed = InputBuffer.MouseHeld
					// reset the mouse held
					InputBuffer.MouseHeld = 0
					InputBuffer.MouseDepth = Screen.CalculateMouseDepth()
				}
			// if we are holding each button then set it to that
			case tcell.Button1:
				InputBuffer.MouseHeld = '1'
			case tcell.Button2:
				InputBuffer.MouseHeld = '2'
			}

			// update the mouse position every mouse event
			x, y := ev.Position()
			InputBuffer.MousePos = V2i(x, y)
			break
		case *tcell.EventResize:
			Screen.Width, Screen.Height = ev.Size()
			Screen.Resize()
			break
		}
		ScreenMutex.Unlock()
	}
}

func (Screen *Screen) Close(){
}