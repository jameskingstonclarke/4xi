package src

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"reflect"
)

type Unit struct {
	*Entity
	*SyncComp
	*PosComp
	*MovementComp
	*HealthComponent
	*AttackComponent
	*RenderComp
}

// used for anything that has health e.g. units, settlements etc
type HealthComponent struct {
	Health float64
}

func (H *HealthComponent) Test(){}

type AttackComponent struct {
	Damage float64
}

func (A *AttackComponent) Test(){}

type UnitSys struct {
	*SystemBase
	SyncComps     []*SyncComp
	PosComps      []*PosComp
	MovementComps []*MovementComp
	HealthComps   []*HealthComponent
	AttackComps   []*AttackComponent
	SelectedUnit  uint32
}

func (ECS *ECS) CreateUnit(pos Vec, dirty bool) *Unit{
	unit := &Unit{
		Entity: ECS.NewEntity("unit"),
		// we hide the RenderComponent from network syncs
		SyncComp: &SyncComp{Dirty: dirty, Hidden: map[string]struct{}{"RenderComp": {}}},
		PosComp: &PosComp{Pos: pos},
		MovementComp: &MovementComp{Target: pos, Speed:  1},
		HealthComponent: &HealthComponent{Health: 1},
		AttackComponent: &AttackComponent{Damage: .25},
		RenderComp: nil,
	}
	return unit
}

func (ECS *ECS) AddUnit(unit *Unit) uint32{
	if ECS.HostMode == CLIENT{
		unit.RenderComp = &RenderComp{Depth:  UNITS_DEPTH, Buffer: FillBufRune('u', tcell.StyleDefault)}
	}
	// register the entity to the ECS
	ECS.AddEntity(unit.Entity, unit.SyncComp, unit.PosComp, unit.MovementComp, unit.HealthComponent, unit.AttackComponent, unit.RenderComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *UnitSys:
			s.AddEntity(unit.Entity, unit.SyncComp, unit.PosComp, unit.MovementComp, unit.HealthComponent, unit.AttackComponent)
		case *NetworkSys:
			s.AddEntity(unit.Entity, unit.SyncComp)
		case *RendererSys:
			s.AddEntity(unit.Entity, unit.RenderComp, unit.PosComp)
		}
	}
	return unit.ID
}

func (U *UnitSys) AddEntity(Entity *Entity, SyncComp *SyncComp, PosComp *PosComp, MovementComp *MovementComp, HealthComp *HealthComponent, AttackComp *AttackComponent){
	U.Entities = append(U.Entities, Entity)
	U.SyncComps = append(U.SyncComps, SyncComp)
	U.PosComps = append(U.PosComps, PosComp)
	U.MovementComps = append(U.MovementComps, MovementComp)
	U.HealthComps = append(U.HealthComps, HealthComp)
	U.AttackComps = append(U.AttackComps, AttackComp)
	U.Size++
}

func (U *UnitSys) Init(){
	// register the unit to the ecs
	U.ECS.RegisterEntity("unit",
		reflect.TypeOf(&Unit{}),
		reflect.ValueOf(&Unit{}).Elem(),
	)
	// initialise the selected unit to nil essentially
	U.SelectedUnit = 1 << 31
	if U.ECS.HostMode == SERVER {
		U.ECS.AddUnit(U.ECS.CreateUnit(V2i(20, 20), true))
		U.ECS.AddUnit(U.ECS.CreateUnit(V2i(10, 10), true))
	}
}

func (U *UnitSys) Update(){}
func (U *UnitSys) Remove(){}

func (U *UnitSys) ListenClickEvent(event ClickEvent){
	if event.Layer == UNITS_DEPTH && event.Type == PRESS && event.Button == '2'{
		for i:=0;i<U.Size;i++{
			// see if we clicked on a unit
			if U.PosComps[i].Pos.Equals(event.WorldPos){
				// select the unit
				if U.SelectedUnit == 1<<31{
					CLog("selected unit ", U.Entities[i].ID)
					U.SelectedUnit = U.Entities[i].ID
				}
			}
		}
	// we are moving the unit
	}else if event.Layer == CELL_DEPTH && event.Type == PRESS && event.Button == '2' && U.SelectedUnit != 1<<31{
		for i:=0; i<U.Size; i++{
			if U.Entities[i].ID == U.SelectedUnit{
				CLog("sending move unit command!")
				id := U.Entities[i].ID
				U.SelectedUnit = 1<<31
				// broadcast a move command to the server
				U.ECS.Event(ClientCommandEvent{
					EventBase: EventBase{},
					Side:      CLIENT,
					Type:      CLIENT_CMD_MOVE_UNIT,
					// id is ID of unit, dest is the position to move to
					Data:      fmt.Sprintf("{\"id\":%d, \"dest\":[%f,%f]}", id, event.WorldPos.X, event.WorldPos.Y),
				})
			}
		}
	}
}

// TODO SERVER
// listen for when the player wants to move a unit
func (U *UnitSys) ListenClientCommandEvent(event ClientCommandEvent){
	if event.Side == SERVER {
		switch event.Type {
		case CLIENT_CMD_MOVE_UNIT:
			SLog("received command to move unit")
			var result map[string]interface{}
			json.Unmarshal([]byte(event.Data), &result)
			id := uint32(result["id"].(float64))
			pos := result["dest"].([]interface{})
			// TODO validate the movement
			for i:=0; i<U.Size; i++{
				if U.Entities[i].ID == id{
					SLog("set entity ", id, " target to ", pos)
					U.MovementComps[i].Target = V2i(int(pos[0].(float64)), int(pos[1].(float64)))
					// TODO for now we just move the unit for demonstration purposes
					U.PosComps[i].Pos = U.MovementComps[i].Target
					// mark the unit as dirty
					U.SyncComps[i].Dirty = true
					// TODO should this be done now? if so we may be triggering quite a lot of syncs
					// trigger a sync event
					U.ECS.Event(ServerCommandEvent{Side:SERVER,Type:SERVER_CMD_SYNC})
				}
			}
		}
	}
}

func (U *UISys) ListenNextTurnEvent(event NextTurnEvent){}