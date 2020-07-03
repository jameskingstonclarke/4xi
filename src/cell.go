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
	Type         uint32
	Arable       float64
	Contaminated float64
}

func (ECS *ECS) AddCell(pos Vec, cellType uint32){
	var style tcell.Style
	var yield rune
	switch cellType{
	case CELL_WATER:
		style = tcell.StyleDefault.Background(tcell.ColorBlue)
		yield = '0'
	case CELL_BEACH:
		style = tcell.StyleDefault.Background(tcell.ColorBeige)
		yield = '1'
	case CELL_PLAINS:
		style = tcell.StyleDefault.Background(tcell.ColorGreen)
		yield = '2'
	}
	cell := &Cell{
		Entity:     NewEntity(),
		PosComp:    &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
			View: WORLD_VIEW,
		},
		RenderComp: &RenderComp{Depth: CELL_DEPTH, Buffer: FillBufRune(yield, style)},//FillBufRune(tcell.RuneBlock, style)},
		CellDatComp: &CellDatComp{Type: cellType},
	}
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *RendererSys:
			s.AddEntity(cell.Entity, cell.RenderComp, cell.PosComp)
		case *WorldSys:
			s.AddEntity(cell.Entity, cell.CellDatComp, cell.PosComp)
		}
	}
}