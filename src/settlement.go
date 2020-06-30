package src

import (
	"github.com/gdamore/tcell"
)

// settlement system
type SettlementSys struct {
	*SystemBase
	PosComps []*PosComp
	SettlementStatsComps []*SettlementStatsComp
}

// empire entity
type Settlement struct{
	*Entity
	*PosComp
	*SettlementStatsComp
	*RenderComp
}


// component for storing empire statistics
type SettlementStatsComp struct {
	Name  string
}

func (ECS *ECS) AddSettlement(name string, pos Vec){

	b := Buf(len(name),2)
	b = BufText(b, name, tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(0,0))
	b = BufRune(b, '▲',tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(len(name)/2, 1))

	settlement := &Settlement{
		Entity:          NewEntity(),
		PosComp:         &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		SettlementStatsComp: &SettlementStatsComp{Name: name,},
		// we have to adjust the position so the rune is the position of the settlement
		RenderComp: &RenderComp{Depth: 0, Pos: pos.Sub(V2i(len(name)/2, 1)), View: WORLD_VIEW, Buffer: b},
	}
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *SettlementSys:
			s.AddEntity(settlement.Entity, settlement.PosComp, settlement.SettlementStatsComp)
		case *RendererSys:
			s.AddEntity(settlement.Entity, settlement.RenderComp)
		}
	}
}

func (S *SettlementSys) Init(){
	S.ECS.AddSettlement("cairo", V2i(0,0))
	S.ECS.AddSettlement("tokyo", V2i(20,15))
}


func (S *SettlementSys) AddEntity(Entity *Entity, PosComp *PosComp, SettlementStatsComp *SettlementStatsComp){
	S.Entities = append(S.Entities, Entity)
	S.PosComps = append(S.PosComps, PosComp)
	S.SettlementStatsComps = append(S.SettlementStatsComps, SettlementStatsComp)
	S.Size++
}

func (S *SettlementSys) Update(){
}

func (S *SettlementSys) Remove(){
}

func (S *SettlementSys) Close(){

}

// listen for sync events
func (S *SettlementSys) Listen(event SyncEvent){}