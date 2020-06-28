package src

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	GameState  *GameState
	Connection net.Conn
	Entities   []Entity
}

func (Client *Client) Init(){
	//conn, err := net.Dial("tcp", "localhost:"+PORT)
	//if err != nil{
	//	LogErr(err)
	//}
	//Client.Connection = conn



	Client.Entities = append(Client.Entities,NewSettlement(nil, "babylon", V2(10,10)))
}

// process all updatable entities
func (Client *Client) Process(){
	for _, entity := range Client.Entities{
		entity.Update()
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