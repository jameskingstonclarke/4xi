package src

import "github.com/gdamore/tcell"

const (
	CELL_WATER  uint32 = 0x0
	CELL_PLAINS uint32 = 0x1
	CELL_BEACH  uint32 = 0x2
)

type Cell struct {
	*Entity
	*PosComp
	*RenderComp
	*CellDatComp
}

type CellDatComp struct {
	Type uint32
}

func (ECS *ECS) AddCell(pos Vec, cellType uint32){
	var style tcell.Style
	switch cellType{
	case CELL_WATER:
		style = tcell.StyleDefault.Foreground(tcell.ColorBlue)
	case CELL_BEACH:
		style = tcell.StyleDefault.Foreground(tcell.ColorBeige)
	case CELL_PLAINS:
		style = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	}
	cell := &Cell{
		Entity:     NewEntity(),
		PosComp:    &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		RenderComp: &RenderComp{Depth: 0, Pos: pos, View: WORLD_VIEW, Buffer: FillBufRune(tcell.RuneBlock, style)},
		CellDatComp: &CellDatComp{Type: cellType},
	}
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *RendererSys:
			s.AddEntity(cell.Entity, cell.RenderComp)
		case *WorldSys:
			s.AddEntity(cell.Entity, cell.CellDatComp)
		}
	}
}