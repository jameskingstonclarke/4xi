package src

import "github.com/gdamore/tcell"

const (
	CELL_WATER  uint32 = 0x0
	CELL_PLAINS uint32 = 0x1
	CELL_BEACH  uint32 = 0x2
)

type Cell struct {
	*Entity
	*SyncComp
	*PosComp
	*RenderComp
	*CellDatComp
}

type CellDatComp struct {
	Type         uint32
	Arable       float64
	Contaminated float64
}

func (ECS *ECS) CreateCell(pos Vec, cellType uint32, dirty bool) *Cell{
	cell := &Cell{
		Entity:     ECS.NewEntity("cell"),
		SyncComp: &SyncComp{Dirty: dirty},
		PosComp:    &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
			View: WORLD_VIEW,
		},
		RenderComp: nil,
		CellDatComp: &CellDatComp{Type: cellType},
	}
	return cell
}

func (ECS *ECS) AddCell(cell *Cell){
	// if we are the client, we want to add a render component
	if ECS.HostMode == CLIENT{
		SLog("client adding cell, ", cell.Pos)
		var yield rune
		var style tcell.Style
		switch cell.CellDatComp.Type{
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
		cell.RenderComp = &RenderComp{Depth: CELL_DEPTH, Buffer: FillBufRune(yield, style)}
	}
	ECS.AddEntity(cell.Entity, cell.SyncComp, cell.PosComp, cell.RenderComp, cell.CellDatComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *NetworkSys:
			s.AddEntity(cell.Entity, cell.SyncComp)
		case *RendererSys:
			s.AddEntity(cell.Entity, cell.RenderComp, cell.PosComp)
		case *WorldSys:
			s.AddEntity(cell.Entity, cell.CellDatComp, cell.PosComp)
		}
	}
}