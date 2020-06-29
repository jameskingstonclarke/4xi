package src

// empire system
type EmpireSys struct {
	*SystemBase
	PosComps []*PosComp
	EmpireStatsComps []*EmpireStatsComp
}

// empire entity
type Empire struct{
	*Entity
	*PosComp
	*EmpireStatsComp
}

func (ECS *ECS) AddEmpire(name string, pos Vec){
	empire := &Empire{
		Entity:          NewEntity(),
		PosComp:         &PosComp{
			Pos: pos,
			Facing: V2i(0,0),
		},
		EmpireStatsComp: &EmpireStatsComp{
			Name: name,
			Money: 0,
		},
	}
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *EmpireSys:
			s.AddEntity(empire.Entity, empire.PosComp, empire.EmpireStatsComp)
		}
	}
}

type GoldenAgeEvent struct {
	EventBase
	// the empire that achieved the golden age
	EmpireStatsComp *EmpireStatsComp
}

// component for storing empire statistics
type EmpireStatsComp struct {
	Name  string
	Money float64
}

func (E *EmpireSys) Init(ECS *ECS){
	ECS.AddEmpire("egypt", V2i(0,0))
}


func (E *EmpireSys) AddEntity(Entity *Entity, PosComp *PosComp, EmpireStatsComp *EmpireStatsComp){
	E.Entities = append(E.Entities, Entity)
	E.PosComps = append(E.PosComps, PosComp)
	E.EmpireStatsComps = append(E.EmpireStatsComps, EmpireStatsComp)
	E.Size++
}

func (E *EmpireSys) Update(ECS *ECS){
	for i:=0;i<E.Size;i++{
		// trigger a golden age on this empire
		if E.EmpireStatsComps[i].Money == 0{
			ECS.Event(GoldenAgeEvent{EventBase{E.Entities[i]}, E.EmpireStatsComps[i]})
		}
	}
}

func (E *EmpireSys) Remove(){
}

func (E *EmpireSys) Close(ECS *ECS){

}
