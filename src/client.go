package src
import "C"

type Client struct {
	ECS *ECS
	Screen *Screen
}

func NewClient(addr string) *Client{
	Screen := &Screen{}
	// create a client by registering all the relevant ECS systems
	ecs := NewECS(CLIENT)
	ecs.RegisterSystem(&NetworkSys{SystemBase: NewSysBase(ecs), ServerAddress: addr})
	ecs.RegisterSystem(&StateSys{SystemBase: NewSysBase(ecs)})
	//ecs.RegisterSystem(&PlayerSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&UnitSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&WorldSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&EmpireSys{SystemBase: NewSysBase(ecs)})
	ecs.RegisterSystem(&SettlementSys{SystemBase: NewSysBase(ecs)})
	
	ecs.RegisterSystem(&RendererSys{SystemBase: NewSysBase(ecs), Screen: Screen})
	ecs.RegisterSystem(&UISys{SystemBase: NewSysBase(ecs), UIManager: NewUIManager(Screen)})
	return &Client{
		ECS:        ecs,
		Screen: Screen,
	}
}

func (Client *Client) Init(){
	Client.Screen.Init()
	go Client.Screen.Poll()
	Client.ECS.Init()
}

// process all updatable entities
func (Client *Client) Process() {
	Client.ECS.Update()
}

func (Client *Client) Close(){
	Client.ECS.Close()
}