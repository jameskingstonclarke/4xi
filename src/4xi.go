package src

const (
	CLIENT = 0x0
	HOST   = 0x1
)

type GameInterface struct {
	Client *Client
	Server *Server
}

func (G *GameInterface) Process(){
	for Running {
		G.Client.Process()
	}
}

func (G *GameInterface) Init(){
	G.Client = NewClient()
	G.Client.Init()
}

func (G *GameInterface) Close(){
	G.Client.Close()
}

var (
	Mode = HOST
	Running = true
)

func Run(){
	InitLogs()
	Log("4xi v_a_001")
	g := &GameInterface{}
	g.Init()
	g.Process()
	g.Close()
	CloseLogs()
}