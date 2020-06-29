package src

import (
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
	Empires []*Empire
	Cells   [][]Cell
}

func NewWorld(width, height int, seed int64) *World{
	w := &World{
		Empires: nil,
		Cells: nil,
	}
	w.Cells= make([][]Cell, width)
	for i := range w.Cells {
		w.Cells[i] = make([]Cell, height)
	}

	// generate the cells
	p := perlin.NewPerlin(2,3,5, seed)
	for x:=0;x<width;x++{
		for y:=0;y<height;y++{
			w.Cells[x][y].Pos = V2(x,y)
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
	return w
}

func (World *World) Update() {}

func (World *World) Draw() {
	for _, cellRow := range World.Cells {
		for _, cell := range cellRow {
			switch cell.Type {
			case CELL_WATER:
				ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorBlue), WORLD_VIEW)
			case CELL_BEACH:
				ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorBeige), WORLD_VIEW)
			case CELL_PLAINS:
				ScreenInstance.Char(tcell.RuneBlock, cell.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW)
			}
		}
	}
}
