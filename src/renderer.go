package src

import (
	"github.com/gdamore/tcell"
)

func Buf(width, height int) *tcell.CellBuffer {
	b := &tcell.CellBuffer{}
	b.Resize(width, height)
	return b
}

// renderer is a system
type RendererSys struct {
	*SystemBase
	RenderComps []*RenderComp
	Screen Screen
}

type RenderComp struct {
	Depth int
	Pos   Vec
	View  uint8
	Buffer *tcell.CellBuffer
}

// fill a cell buffer with text
func BufText(buf *tcell.CellBuffer, text string, style tcell.Style, pos Vec) *tcell.CellBuffer{
	for i, r := range text{
		buf.SetContent(int(pos.X)+i, int(pos.Y), r, nil, style)
	}
	return buf
}

// fill a cell buffer with text
func BufRune(buf *tcell.CellBuffer, rune rune, style tcell.Style, pos Vec) *tcell.CellBuffer{
	buf.SetContent(int(pos.X), int(pos.Y), rune, nil, style)
	return buf
}

// fill a cell buffer with text
func FillBufRune(rune rune, style tcell.Style)*tcell.CellBuffer{
	buf := &tcell.CellBuffer{}
	buf.Resize(1,1)
	buf.SetContent(0, 0, rune, nil, style)
	return buf
}

func (R *RendererSys) Init(){
	screen, err := tcell.NewScreen()
	if err != nil{
		CLogErr(err)
	}
	R.Screen.Screen = screen
	if err = screen.Init(); err != nil{
		CLogErr(err)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	screen.EnableMouse()
	screen.Clear()

	R.Screen.Width = DEFAULT_WIDTH
	R.Screen.Height = DEFAULT_HEIGHT
	R.Screen.Resize()
	R.Screen.ZBuffer = make([]tcell.CellBuffer, 3)
	for i:=0; i<Z_DEPTH; i++ {
		R.Screen.ZBuffer[i] = tcell.CellBuffer{}
		R.Screen.ZBuffer[i].Resize(R.Screen.Width, R.Screen.Height)
	}

	go R.Screen.Poll()
}

func (R *RendererSys) Close(){

}

func (R *RendererSys) AddEntity(Entity *Entity, RenderComp *RenderComp){
	R.Entities = append(R.Entities, Entity)
	R.RenderComps = append(R.RenderComps, RenderComp)
	R.Size++
}

func (R *RendererSys) Update(){

	// process camera movement
	if InputBuffer.KeyPressed == 'a'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(1,0))
	}else if InputBuffer.KeyPressed == 'd'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(-1,0))
	}else if InputBuffer.KeyPressed == 'w'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(0,1))
	}else if InputBuffer.KeyPressed == 's'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(0,-1))
	}else if InputBuffer.CtrlKeyPressed == tcell.KeyEscape{
		Running = false
	}

	// render each RenderComp
	for _, r := range R.RenderComps{
		width, height := r.Buffer.Size()
		// iterate over each cell in the cell buffer
		for x:=0; x<width;x ++{
			for y:=0; y<height;y ++{
				rune, _, style, _ := r.Buffer.GetContent(x,y)
				// draw at the offset from the RenderComp position
				R.Screen.Char(rune, r.Pos.Add(V2i(x,y)), style, r.View, r.Depth)
			}
		}
	}
	R.Screen.Draw()
}

func (R *RendererSys) Remove(){
}