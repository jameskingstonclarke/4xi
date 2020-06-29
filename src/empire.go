package src

// empire system
type EmpireSys struct {
	*SystemBase
}

// empire entity
type Empire struct{
	*Entity
	*PosComp
}

func (E *EmpireSys) Close(ECS *ECS){

}

func (E *EmpireSys) AddEntity(Entity *Entity, PosComp *PosComp){
}

func (E *EmpireSys) Update(ECS *ECS){
}

func (E *EmpireSys) Remove(){
}