package src

import (
	"github.com/aquilax/go-perlin"
)

type WorldSys struct {
	*SystemBase
	CellDatComps []*CellDatComp
	Width, Height int
}
func (W *WorldSys) AddEntity(Entity *Entity, CellDatComp *CellDatComp){
	W.Entities = append(W.Entities, Entity)
	W.CellDatComps = append(W.CellDatComps, CellDatComp)
	W.Size++
}

func (W *WorldSys) Init(ECS *ECS){
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
			ECS.AddCell(V2i(x,y), cellType)
		}
	}
}

func (W *WorldSys) Update(ECS *ECS){
	if InputBuffer.KeyPressed == 'u'{
		// update here
	}
}

func (W *WorldSys) Remove(){
}

// low priority as we want to render last
func (W *WorldSys) Priority() int {
	return 0
}