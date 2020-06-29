package src


var(
	Running = true
)

type Game struct {
	Client *Client
	Server *Server
}

func (G *Game) Process(){
	for Running {
		G.Client.Process()
	}
}

func (G *Game) Init(){
	G.Client = NewClient()
	G.Client.Init()
}

func (G *Game) Close(){
	G.Client.Close()
}