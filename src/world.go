package src

import (
	"fmt"
	"github.com/aquilax/go-perlin"
	"reflect"
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
	W.ECS.RegisterEntity("cell", reflect.TypeOf(&Cell{}), reflect.ValueOf(&Cell{}).Elem())
	if W.ECS.HostMode == SERVER {
		W.Width = 50
		W.Height = 50
		// generate the cells
		p := perlin.NewPerlin(2, 5, 5, 1)
		for x := 0; x < W.Width; x++ {
			for y := 0; y < W.Height; y++ {
				var cellType uint32
				noise := p.Noise2D(float64(x)/150, float64(y)/150) * -1
				if noise > 0.2 {
					cellType = CELL_WATER
				} else if noise > 0.1 {
					cellType = CELL_BEACH
				} else {
					cellType = CELL_PLAINS
				}
				W.ECS.AddCell(W.ECS.CreateCell(V2i(x, y), cellType, true))
			}
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
	if event.Layer == CELL_DEPTH && event.Type == PRESS && event.Button =='1'{
		for i := 0; i < W.Size; i++ {
			if event.WorldPos.Equals(W.PosComps[i].Pos) {
				W.ECS.Event(NewWinEvent{
					ID:    fmt.Sprintf("cell: %f, %f", event.WorldPos.X, event.WorldPos.Y),
					Title: fmt.Sprintf("cell: %f, %f", event.WorldPos.X, event.WorldPos.Y),
					Text: map[interface{}]func(){
						fmt.Sprintf("type %d", W.CellDatComps[i].Type):nil,
						fmt.Sprintf("arable %f", W.CellDatComps[i].Arable):nil,
						fmt.Sprintf("contaminated %f", W.CellDatComps[i].Contaminated):nil,
					},
				})
			}
		}
	}
}