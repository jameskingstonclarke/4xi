package src

import (
	"fmt"
	"github.com/gdamore/tcell"
)

// settlement system
type SettlementSys struct {
	*SystemBase
	PosComps []*PosComp
	SettlementStatsComps []*SettlementStatsComp
	SelectedCity uint32
}

// empire entity
type Settlement struct{
	*Entity
	*SyncComp
	*PosComp
	*SettlementStatsComp
	*RenderComp
}

// component for storing empire statistics
type SettlementStatsComp struct {
	Name       string
	Population float64
	Production float64
}

func (ECS *ECS) AddSettlement(name string, pos Vec){

	b := Buf(len(name),2)
	b = BufText(b, name, tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(0,0))
	b = BufRune(b, '▲',tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(len(name)/2, 1))

	settlement := &Settlement{
		Entity:          NewEntity(),
		SyncComp: &SyncComp{Dirty: false},
		PosComp:         &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		SettlementStatsComp: &SettlementStatsComp{Name: name},
		RenderComp: &RenderComp{Depth: STRUCTURES_DEPTH, Buffer: b, Offset: V2i(-len(name)/2,-1)},
	}

	ECS.AddEntity(settlement.Entity, settlement.SyncComp, settlement.PosComp, settlement.SettlementStatsComp)

	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *RendererSys:
			s.AddEntity(settlement.Entity, settlement.RenderComp, settlement.PosComp)
		case *NetworkSys:
			s.AddEntity(settlement.Entity, settlement.SyncComp)
		case *SettlementSys:
			s.AddEntity(settlement.Entity, settlement.PosComp, settlement.SettlementStatsComp)
		}
	}
}

func (S *SettlementSys) Init(){
	S.ECS.AddSettlement("cairo", V2i(0,0))
	S.ECS.AddSettlement("tokyo", V2i(20,15))
	S.ECS.AddSettlement("london", V2i(5,12))
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
func (S *SettlementSys) ListenSyncEvent(event SyncEvent){}

func (S *SettlementSys) ListenClickEvent(event ClickEvent){
	// this is the first way of checking the click (this assumes the structures depth is constant)
	if event.Layer == STRUCTURES_DEPTH && event.Type == PRESS{
		for i := 0; i < S.Size; i++ {
			if event.WorldPos.Equals(S.PosComps[i].Pos) {
				// open UI
				if event.Button == '1' {
					S.ECS.Event(NewWinEvent{
						ID:    fmt.Sprintf("settlement: %f, %f", event.WorldPos.X, event.WorldPos.Y),
						Title: fmt.Sprintf("settlement: %f, %f", event.WorldPos.X, event.WorldPos.Y),
						Text: []string{
							fmt.Sprintf("population %f", S.SettlementStatsComps[i].Population),
							fmt.Sprintf("production %f", S.SettlementStatsComps[i].Production),
						},
					})
				}
			}
		}
	}
}