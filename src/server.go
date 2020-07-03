package src

type Server struct {
	ECS        *ECS
}

func NewServer() *Server{
	// create a client by registering all the relevant ECS systems
	ecs := NewECS(SERVER)
	ecs.RegisterSystem(&NetworkSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&StateSys{SystemBase: NewSysBase(ecs)})
	//ecs.RegisterSystem(&PlayerSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&UnitSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&WorldSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&EmpireSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&SettlementSys{SystemBase: NewSysBase(ecs)})
	return &Server{
		ECS:        ecs,
	}
}

func (S *Server) Init(){
	S.ECS.Init()
}

// process all update-able entities
func (S *Server) Process(){
	S.ECS.Update()
}

func (S *Server) Close(){
	S.ECS.Close()
}