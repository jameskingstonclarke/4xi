package src

type Logic struct {
	Client *Client
	Server *Server
}

var (
	LogicInstance = &Logic{}
)

func (Logic *Logic) Process(){
	for Running{
		// lock the screen, and then process all logic
		ScreenMutex.Lock()

		s := &Settlement{
			Empire:     nil,
			Name:       "babylon",
			Population: 1,
			X:          10,
			Y:          10,
		}
		s.Draw()

		// finally release the screen for the game to render
		ScreenMutex.Unlock()
	}
}

func (Logic *Logic) Init(){
	Logic.Client = &Client{GameState: nil}
	// if we are hosting a server, setup the server
	if Mode == HOST {
		Logic.Server = &Server{
			Players:   nil,
			GameState: nil,
		}
	}
}

func (Logic *Logic) Close(){
	WaitGroup.Done()
}