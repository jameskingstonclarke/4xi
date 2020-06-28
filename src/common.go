package src

const (
	PORT = "7777"
)

// represents a player in the game
type Player struct {
	Name string
}

// represents the state of the game, shared by the client & the server
type GameState struct {
	Turn int
}