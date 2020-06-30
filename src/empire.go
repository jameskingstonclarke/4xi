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

func (E *EmpireSys) Init(){
	E.ECS.AddEmpire("egypt", V2i(0,0))
}


func (E *EmpireSys) AddEntity(Entity *Entity, PosComp *PosComp, EmpireStatsComp *EmpireStatsComp){
	E.Entities = append(E.Entities, Entity)
	E.PosComps = append(E.PosComps, PosComp)
	E.EmpireStatsComps = append(E.EmpireStatsComps, EmpireStatsComp)
	E.Size++
}

func (E *EmpireSys) Update(){
	for i:=0;i<E.Size;i++{
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