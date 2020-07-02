package src

// all ECS code related to spacial stuff e.g. movements, positions etc

type PosComp struct {
	Pos    Vec
	Facing Vec
	// world or screen
	View   uint8
}

type MovementComp struct{
	Target Vec
	Speed  float64
}

type MovementSys struct {
	*SystemBase
	// store all the position & movement components components
	PosComps 	  []*PosComp
	MovementComps []*MovementComp
}

func (MovementSys *MovementSys) Add(Entity *Entity, PosComp *PosComp, MovementComp *MovementComp){
	MovementSys.Entities = append(MovementSys.Entities, Entity)
	MovementSys.PosComps = append(MovementSys.PosComps, PosComp)
	MovementSys.MovementComps = append(MovementSys.MovementComps, MovementComp)
	MovementSys.Size++
}

// TODO fix this, we only want to apply the movement on a next turn event
func (MovementSys *MovementSys) Update(){
	// iterate over each entity
	for i:=0;i<MovementSys.Size;i++{
		// get the direction between the 2 vectors
		dir := MovementSys.MovementComps[i].Target.Sub(MovementSys.PosComps[i].Pos)
		dir = dir.Normalize().Round()
		MovementSys.PosComps[i].Pos = MovementSys.PosComps[i].Pos.Add(dir)
	}
}

func (MovementSys *MovementSys) Remove(ECS *ECS){}