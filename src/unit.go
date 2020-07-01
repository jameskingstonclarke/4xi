package src

type Unit struct {
	*Entity
	*SyncComp
	*PosComp
	*MovementComp
	*RenderComp
}

func (ECS *ECS) AddUnit(){
	unit := &Unit{
		Entity:          NewEntity(),
	}
	// register the entity to the ECS
	ECS.AddEntity(unit.Entity, unit.PosComp, unit.MovementComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *NetworkSys:
			s.AddEntity(unit.Entity, unit.SyncComp)
		}
	}
}

