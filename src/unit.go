package src

import "github.com/gdamore/tcell"

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

type AttackComponent struct {
	Damage float64
}

type UnitSys struct {
	*SystemBase
	PosComps      []*PosComp
	MovementComps []*MovementComp
	HealthComps   []*HealthComponent
	AttackComps   []*AttackComponent
	SelectedUnit  uint32
}

func (ECS *ECS) AddUnit(pos Vec) uint32{
	unit := &Unit{
		Entity:          NewEntity(),
		PosComp: &PosComp{Pos: pos},
		MovementComp: &MovementComp{Target: pos, Speed:  1},
		HealthComponent: &HealthComponent{Health: 1},
		AttackComponent: &AttackComponent{Damage: .25},
		RenderComp: &RenderComp{Depth:  UNITS_DEPTH, Buffer: FillBufRune('u', tcell.StyleDefault)},
	}
	// register the entity to the ECS
	ECS.AddEntity(unit.Entity, unit.PosComp, unit.MovementComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *UnitSys:
			s.AddEntity(unit.Entity, unit.PosComp, unit.MovementComp, unit.HealthComponent, unit.AttackComponent)
		case *NetworkSys:
			s.AddEntity(unit.Entity, unit.SyncComp)
		case *RendererSys:
			s.AddEntity(unit.Entity, unit.RenderComp, unit.PosComp)
		}
	}
	return unit.ID
}

func (U *UnitSys) AddEntity(Entity *Entity, PosComp *PosComp, MovementComp *MovementComp, HealthComp *HealthComponent, AttackComp *AttackComponent){
	U.Entities = append(U.Entities, Entity)
	U.PosComps = append(U.PosComps, PosComp)
	U.MovementComps = append(U.MovementComps, MovementComp)
	U.HealthComps = append(U.HealthComps, HealthComp)
	U.AttackComps = append(U.AttackComps, AttackComp)
	U.Size++
}

func (U *UnitSys) Init(){
	// initialise the selected unit to nil essentially
	U.SelectedUnit = 1<<31
	U.ECS.AddUnit(V2i(20,20))
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
				// move the unit
				U.PosComps[i].Pos = event.WorldPos
				U.SelectedUnit = 1<<31
			}
		}
	}
}