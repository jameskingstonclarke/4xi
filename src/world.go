package src

import (
	"fmt"
	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell"
)

const (
	CELL_WATER  = 0x0
	CELL_PLAINS = 0x1
	CELL_BEACH = 0x2
)

type Cell struct {
	Pos  Vec
	Type uint8
}

type World struct {
	Empires 	  []*Empire
	Width, Height int
	Cells   	  [][]Cell
	DrawManager   *DrawManager
}

func NewWorld(width, height int, seed int64) *World{
	w := &World{
		Empires: nil,
		Width: width,
		Height: height,
		Cells: nil,
		DrawManager: NewDrawManager(),
	}
	w.Cells= make([][]Cell, width)
	for i := range w.Cells {
		w.Cells[i] = make([]Cell, height)
	}

	// generate the cells
	p := perlin.NewPerlin(2,3,5, seed)
	for x:=0;x<width;x++{
		for y:=0;y<height;y++{
			w.Cells[x][y].Pos = V2i(x,y)
			noise := p.Noise2D(float64(x)/100,float64(y)/100)*-1
			if noise > 0.2{
				w.Cells[x][y].Type = CELL_WATER
			}else if noise > 0.1{
				w.Cells[x][y].Type = CELL_BEACH
			}else{
				w.Cells[x][y].Type = CELL_PLAINS
			}
		}
	}
	w.InitUI()
	return w
}

func (W *World) Update() {
	for _, e := range W.Empires{
		e.Update()
	}
}

func (W *World) InitUI(){
	// add the window for clicking on a tile
	w:=W.DrawManager.NewWin("cell_inspector", false, V2i(10, 10), V2i(20, 20), SCREEN_VIEW)
	w.NewText("pos", "pos", true, V2(0, 0), SCREEN_VIEW, tcell.StyleDefault, nil)

	// Init the UI for each world cell
	for x:=0;x<W.Width;x++{
		for y:=0;y<W.Height;y++{
			cell := W.Cells[x][y]

			var style tcell.Style
			switch cell.Type{
			case CELL_WATER:
				style = tcell.StyleDefault.Foreground(tcell.ColorBlue)
			case CELL_BEACH:
				style = tcell.StyleDefault.Foreground(tcell.ColorBeige)
			case CELL_PLAINS:
				style = tcell.StyleDefault.Foreground(tcell.ColorGreen)
			}

			W.DrawManager.NewText(fmt.Sprintf("%f",cell.Pos), string(tcell.RuneBlock), true, cell.Pos, WORLD_VIEW, style,  func() {
				w.Enable(true)
			})
		}
	}
}


func (W *World) Draw() {
	//for _, cellRow := range World.Cells {
	//	//	for _, cell := range cellRow {
	//	//		switch cell.Type {
	//	//		case CELL_WATER:
	//	//			ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorBlue), WORLD_VIEW)
	//	//		case CELL_BEACH:
	//	//			ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorBeige), WORLD_VIEW)
	//	//		case CELL_PLAINS:
	//	//			ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW)
	//	//		}
	//	//	}
	//	//}
	W.DrawManager.Draw()
	for _, empire := range W.Empires{
		empire.Draw()
	}
}
