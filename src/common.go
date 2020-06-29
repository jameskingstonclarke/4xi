package src

const (
	PORT = "7777"

	SERVER_MODE = 0x0
	CLIENT_MODE = 0x1
)

// represents a player in the game
type Player struct {
	Name string
}

// represents the state of the game, shared by the client & the server
type GameState struct {
	Mode uint8
	Turn int
	// represents every single entity in the game (e.g. settlements, empires etc)
	World *World
}

func (GameState *GameState) Update(){
	GameState.NextTurn()
}

func (GameState *GameState) NextTurn(){
	GameState.Turn++
	GameState.World.Update()
}

func (GameState *GameState) Draw(){
	GameState.World.Draw()
}