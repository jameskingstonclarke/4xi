package src

import (
	"github.com/gdamore/tcell"
	"os"
	"sync"
)

const (
	DEFAULT_WIDTH  = 200
	DEFAULT_HEIGHT = 100

	WORLD_VIEW     = 0x0
	SCREEN_VIEW    = 0x1
)

type InputData struct {
	MousePressed   tcell.ButtonMask
	KeyPressed     rune
	// used for special key presses e.g. ctrl + c
	CtrlKeyPressed tcell.Key
	MousePos       Vec
}

type Screen struct {
	Screen        tcell.Screen
	Cam    		  Vec
	Width, Height int
	CellBuffer    tcell.CellBuffer
	InputBuffer   InputData
}

var (
	ScreenInstance = &Screen{}
	ScreenMutex = &sync.Mutex{}
)

func (Screen *Screen) Init(){
	screen, err := tcell.NewScreen()
	if err != nil{
		LogErr(err)
	}
	Screen.Screen = screen
	if err = screen.Init(); err != nil{
		LogErr(err)
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
	Screen.CellBuffer = tcell.CellBuffer{}
	Screen.CellBuffer.Resize(Screen.Width, Screen.Height)
}

func (Screen *Screen) Resize(){
	Screen.Screen.Sync()
}

func (Screen *Screen) WorldToScreen(v Vec) Vec{
	vNew := v.Add(Screen.Cam)
	return vNew.Add(Vec{X:Screen.Width/2, Y:Screen.Height/2})
}

func (Screen *Screen) ScreenToWorld(v Vec) Vec{
	vNew := v.Sub(Vec{X:Screen.Width/2, Y:Screen.Height/2})
	return vNew.Sub(Screen.Cam)
}

// TODO We may be able to just write the text directly
func (Screen *Screen) Text(text string, pos Vec, style tcell.Style, view uint8){
	for i, r := range text {
		if pos.X+i >= Screen.Width {
			pos.X=0
			pos.Y++
		}
		Screen.Char(r, V2(pos.X+i, pos.Y), style, view)
	}
}

func (Screen *Screen) Char(r rune, pos Vec, style tcell.Style, view uint8){
	if view == WORLD_VIEW{
		pos = Screen.WorldToScreen(pos)
	}
	Screen.CellBuffer.SetContent(pos.X, pos.Y, r, nil, style)
}

func (Screen *Screen) Rect(r rune, pos Vec, width, height int, style tcell.Style, fill bool, view uint8){
	for xTmp:=pos.X; xTmp<pos.Y+width; xTmp++ {
		for yTmp:=pos.X; yTmp<pos.Y+height; yTmp++ {
			if fill {
				Screen.Char(r, V2(xTmp, yTmp), style, view)
			}else{
				if (xTmp==pos.X || xTmp==pos.X+width-1) || (yTmp==pos.Y || yTmp==pos.Y+height-1) {
					Screen.Char(r, V2(xTmp, yTmp), style, view)
				}
			}
		}
	}
}

func (Screen *Screen) Draw(){
	ScreenInstance.InputBuffer = InputData{
		MousePressed:   0,
		KeyPressed:     0,
		CtrlKeyPressed: 0,
		MousePos:       ScreenInstance.InputBuffer.MousePos,
	}
	for y := 0; y < Screen.Height; y++ {
		for x := 0; x < Screen.Width; x++ {
			rune, _, style, _ := Screen.CellBuffer.GetContent(x,y)
			Screen.Screen.SetCell(x, y, style, rune)
		}
	}
	Screen.Screen.Show()
	//Screen.Screen.Clear()
	Screen.CellBuffer = tcell.CellBuffer{}
	Screen.CellBuffer.Resize(Screen.Width, Screen.Height)
}

func (Screen *Screen) Poll() {
	defer WaitGroup.Done()
	for Running {
		ev := Screen.Screen.PollEvent()
		switch ev := ev.(type){
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune{
				ScreenInstance.InputBuffer.KeyPressed = ev.Rune()
			}else if ev.Key() == tcell.KeyEscape{
				os.Exit(2)
			}else {
				ScreenInstance.InputBuffer.CtrlKeyPressed = ev.Key()
			}
			break
		case *tcell.EventMouse:
			ScreenInstance.InputBuffer.MousePressed = ev.Buttons()
			x, y := ev.Position()
			ScreenInstance.InputBuffer.MousePos = V2(x, y)
			break
		case *tcell.EventResize:
			Screen.Width, Screen.Height = ev.Size()
			Screen.Resize()
			break
		}
	}
}

func (Screen *Screen) Close(){
	WaitGroup.Done()
}