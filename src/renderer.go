package src

import (
	"github.com/gdamore/tcell"
)

const (
	PRESS = 0x0
	HELD  = 0x1

	CAM_SPEED = 0x3

	CELL_DEPTH 		 = 0x0
	STRUCTURES_DEPTH = 0x1
	UNITS_DEPTH 	 = 0x2
	UI_DEPTH    	 = 0x3
)

func Buf(width, height int) *tcell.CellBuffer {
	b := &tcell.CellBuffer{}
	b.Resize(width, height)
	return b
}

type ClickEvent struct{
	EventBase
	Button      rune
	ScreenPos   Vec
	WorldPos    Vec
	Layer 	    int
	Type  	    uint8
}


// renderer is a system
type RendererSys struct {
	*SystemBase
	PosComps    []*PosComp
	RenderComps []*RenderComp
	Screen *Screen
	Clicked tcell.ButtonMask
}

type RenderComp struct {
	Depth int
	Offset  Vec
	Buffer *tcell.CellBuffer
}

func (R *RenderComp) Deserialize(data interface{}){}

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

func (R *RendererSys) Init(){}
func (R *RendererSys) Close(){}

func (R *RendererSys) AddEntity(Entity *Entity, RenderComp *RenderComp, PosComp *PosComp){
	R.Entities = append(R.Entities, Entity)
	R.RenderComps = append(R.RenderComps, RenderComp)
	R.PosComps = append(R.PosComps, PosComp)
	R.Size++
}

func (R *RendererSys) Update(){

	// process camera movement
	if InputBuffer.KeyHeld == 'a'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(CAM_SPEED,0))
	}else if InputBuffer.KeyHeld == 'd'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(-CAM_SPEED,0))
	}else if InputBuffer.KeyHeld == 'w'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(0,CAM_SPEED))
	}else if InputBuffer.KeyHeld == 's'{
		R.Screen.Cam = R.Screen.Cam.Add(V2(0,-CAM_SPEED))
	}else if InputBuffer.CtrlKeyPressed == tcell.KeyEscape{
		Running = false
	}

	// first render each RenderComp to the screen
	for i:=0; i<R.Size; i++{
		r := R.RenderComps[i]
		width, height := r.Buffer.Size()
		// iterate over each cell in the cell buffer
		for x:=0; x<width;x ++{
			for y:=0; y<height;y ++{
				rune, _, style, _ := r.Buffer.GetContent(x,y)
				// draw at the offset from the RenderComp position
				R.Screen.Char(rune, R.PosComps[i].Pos.Add(V2i(x,y)).Add(R.RenderComps[i].Offset), style, R.PosComps[i].View, r.Depth)
			}
		}
	}

	// check for mouse inputs
	if InputBuffer.MousePressed != 0 {
		R.ECS.Event(ClickEvent{
			Button:	   InputBuffer.MousePressed,
			ScreenPos: InputBuffer.MousePos,
			WorldPos:  R.Screen.ScreenToWorld(InputBuffer.MousePos),
			Layer:     InputBuffer.MouseDepth,
		})
	}

	R.Screen.Draw()
}

func (R *RendererSys) Remove(){
}