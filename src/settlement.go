package src

import (
	"fmt"
	"github.com/gdamore/tcell"
	"reflect"
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

func (S *SettlementStatsComp) Test(){}

func (ECS *ECS) CreateSettlement(name string, pos Vec, dirty bool) *Settlement{
	settlement := &Settlement{
		Entity: ECS.NewEntity("settlement"),
		SyncComp: &SyncComp{Dirty: dirty, Hidden: map[string]struct{}{"RenderComp": {}}},
		PosComp:  &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		SettlementStatsComp: &SettlementStatsComp{Name: name},
		RenderComp: nil,
	}
	return settlement
}

func (ECS *ECS) AddSettlement(settlement *Settlement) uint32{

	if ECS.HostMode == CLIENT{
		b := Buf(len(settlement.Name),2)
		b = BufText(b, settlement.Name, tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(0,0))
		b = BufRune(b, '▲',tcell.StyleDefault.Foreground(tcell.ColorRed), V2i(len(settlement.Name)/2, 1))
		settlement.RenderComp = &RenderComp{Depth: STRUCTURES_DEPTH, Buffer: b, Offset: V2i(-len(settlement.Name)/2,-1)}
	}

	ECS.AddEntity(settlement.Entity, settlement.SyncComp, settlement.PosComp, settlement.SettlementStatsComp, settlement.RenderComp)

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
	return settlement.ID
}

func (S *SettlementSys) Init(){
	S.ECS.RegisterEntity("settlement",
		reflect.TypeOf(&Settlement{}),
		reflect.ValueOf(&Settlement{}).Elem())
	if S.ECS.HostMode == SERVER {
		S.ECS.AddSettlement(S.ECS.CreateSettlement("cairo", V2i(0, 0), true))
		S.ECS.AddSettlement(S.ECS.CreateSettlement("tokyo", V2i(20, 15), true))
		S.ECS.AddSettlement(S.ECS.CreateSettlement("london", V2i(5, 12), true))
	}
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
						Text: map[interface{}]func(){
							fmt.Sprintf("population %f", S.SettlementStatsComps[i].Population):nil,
							fmt.Sprintf("production %f", S.SettlementStatsComps[i].Production):nil,
						},
					})
				}
			}
		}
	}
}