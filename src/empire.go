package src

import "reflect"

// empire system
type EmpireSys struct {
	*SystemBase
	PosComps []*PosComp
	EmpireStatsComps []*EmpireStatsComp
}

// empire entity
type Empire struct{
	*Entity
	*SyncComp
	*PosComp
	*EmpireStatsComp
}

// triggered when something should change the stats of an empire
type EmpireStatsEvent struct {
	EventBase
	TargetEmpire   string
	MoneyPerTurnDT float64
}

// component for storing empire statistics
type EmpireStatsComp struct {
	Name  string
	Money float64
}

func (E *EmpireStatsComp) Deserialize(data interface{}){}

func (ECS *ECS) AddEmpire(name string, pos Vec, dirty bool){
	empire := &Empire{
		Entity: ECS.NewEntity("empire"),
		SyncComp: &SyncComp{Dirty: dirty},
		PosComp: &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		EmpireStatsComp: &EmpireStatsComp{
			Name: name,
			Money: 0,
		},
	}
	ECS.AddEntity(empire.Entity, empire.SyncComp, empire.PosComp, empire.EmpireStatsComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *NetworkSys:
			s.AddEntity(empire.Entity, empire.SyncComp)
		case *EmpireSys:
			s.AddEntity(empire.Entity, empire.PosComp, empire.EmpireStatsComp)
		}
	}
}

func (E *EmpireSys) Init(){
	E.ECS.RegisterEntity("empire", reflect.TypeOf(&Empire{}), reflect.ValueOf(&Empire{}).Elem())
	if E.ECS.HostMode == SERVER {
		E.ECS.AddEmpire("egypt", V2i(0, 0), true)
	}
}


func (E *EmpireSys) AddEntity(Entity *Entity, PosComp *PosComp, EmpireStatsComp *EmpireStatsComp){
	E.Entities = append(E.Entities, Entity)
	E.PosComps = append(E.PosComps, PosComp)
	E.EmpireStatsComps = append(E.EmpireStatsComps, EmpireStatsComp)
	E.Size++
}

func (E *EmpireSys) Update(){
	for i:=0;i<int(E.Size);i++{
	}
}

func (E *EmpireSys) Remove(){
}

func (E *EmpireSys) Close(){

}

// listen out for when something modifies the empires stats
func (E *EmpireSys) Listen(Event EmpireStatsEvent){
	for _, e := range E.EmpireStatsComps{
		if e.Name == Event.TargetEmpire{
			// update the empire
			e.Money += Event.MoneyPerTurnDT
		}
	}
}