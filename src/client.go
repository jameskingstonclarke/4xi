package src

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	GameState  *GameState
	Connection net.Conn
}

func (Client *Client) Init(){
	//conn, err := net.Dial("tcp", "localhost:"+PORT)
	//if err != nil{
	//	LogErr(err)
	//}
	//Client.Connection = conn

	Client.GameState = &GameState{
		Turn:     0,
		Entities: nil,
	}

	Client.GameState.Entities = append(Client.GameState.Entities, NewWorld(100,50,123))
	Client.GameState.Entities = append(Client.GameState.Entities, NewSettlement(nil, "babylon", V2(10,10)))
}

// process all updatable entities
func (Client *Client) Process(){

	// process camera movement
	if ScreenInstance.InputBuffer.KeyPressed == 'a'{
		ScreenInstance.Cam = ScreenInstance.Cam.Add(V2(1,0))
	}else if ScreenInstance.InputBuffer.KeyPressed == 'd'{
		ScreenInstance.Cam = ScreenInstance.Cam.Add(V2(-1,0))
	}else if ScreenInstance.InputBuffer.KeyPressed == 'w'{
		ScreenInstance.Cam = ScreenInstance.Cam.Add(V2(0,1))
	}else if ScreenInstance.InputBuffer.KeyPressed == 's'{
		ScreenInstance.Cam = ScreenInstance.Cam.Add(V2(0,-1))
	}

	for _, entity := range Client.GameState.Entities{
		// as we are the client, we only draw the entities
		// we continually wait for the server to send us the updated game state
		entity.Draw()
	}
}

func (Client *Client) SendMsg(msg string){
	fmt.Fprintf(Client.Connection, msg+"\n")
}

// listen from a message from the server
func (Client *Client) ListenMsg() string {
	msg, err := bufio.NewReader(Client.Connection).ReadString('\n')
	if err != nil {
		LogErr(err)
	}
	return msg
}

// request the state of the game from the server
// involves requesting the current turn, the map status etc
func (Client *Client) RequestGameState(){

}

func (Client *Client) Close(){
	Client.Connection.Close()
}