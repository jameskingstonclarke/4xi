package src

type Logic struct {
	Client *Client
	Server *Server
}

var (
	LogicInstance = &Logic{}
)

func (Logic *Logic) Process(){
	for Running {
		Logic.Client.Process()
		ScreenInstance.Draw()
	}
}

func (Logic *Logic) Init(){
	Logic.Client = &Client{GameState: nil}
	Logic.Client.Init()
	// if we are hosting a server, setup the server
	if Mode == HOST {
		Logic.Server = &Server{
			Players:   nil,
			GameState: nil,
		}
		Logic.Server.Init()
	}
}

func (Logic *Logic) Close(){
	Logic.Client.Close()
}