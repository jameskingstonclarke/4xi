package src

import (
	"fmt"
	"github.com/aquilax/go-perlin"
)

type WorldSys struct {
	*SystemBase
	CellDatComps []*CellDatComp
	// position of cells
	PosComps     []*PosComp
	Width, Height int
}
func (W *WorldSys) AddEntity(Entity *Entity, CellDatComp *CellDatComp, PosComp *PosComp){
	W.Entities = append(W.Entities, Entity)
	W.CellDatComps = append(W.CellDatComps, CellDatComp)
	W.PosComps = append(W.PosComps, PosComp)
	W.Size++
}

func (W *WorldSys) Init(){
	W.Width = 200
	W.Height = 100
	// generate the cells
	p := perlin.NewPerlin(2,3,5, 12345)
	for x:=0;x<W.Width;x++{
		for y:=0;y<W.Height;y++{
			var cellType uint32
			noise := p.Noise2D(float64(x)/100,float64(y)/100)*-1
			if noise > 0.2{
				cellType = CELL_WATER
			}else if noise > 0.1{
				cellType = CELL_BEACH
			}else{
				cellType = CELL_PLAINS
			}
			W.ECS.AddCell(V2i(x,y), cellType)
		}
	}
}

func (W *WorldSys) Update(){
	// process each cell here
}

func (W *WorldSys) Remove(){
}

// low priority as we want to render last
func (W *WorldSys) Priority() int {
	return 0
}

func (W *WorldSys) ListenClickEvent(event ClickEvent){
	if event.Layer == CELL_DEPTH && event.Type == PRESS {
		for i := 0; i < W.Size; i++ {
			if event.WorldPos.Equals(W.PosComps[i].Pos) {
				W.ECS.Event(NewWinEvent{
					ID:    fmt.Sprintf("cell: %f, %f", event.WorldPos.X, event.WorldPos.Y),
					Title: fmt.Sprintf("cell: %f, %f", event.WorldPos.X, event.WorldPos.Y),
					Text: []string{fmt.Sprintf("cell type %d", W.CellDatComps[i].Type)},
				})
			}
		}
	}
}