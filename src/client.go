package src

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	ECS        *ECS
	Connection net.Conn
}


func NewClient() *Client{
	//conn, err := net.Dial("tcp", "localhost:"+PORT)
	//if err != nil{
	//	LogErr(err)
	//}
	//Client.Connection = conn

	// create a client by registering all the relevant ECS systems
	ecs := NewECS()
	ecs.RegisterSystem(&RendererSys{SystemBase: NewSysBase()})
	ecs.RegisterSystem(&WorldSys{SystemBase: NewSysBase()})
	ecs.RegisterSystem(&EmpireSys{SystemBase: NewSysBase()})
	ecs.RegisterSystem(&SettlementSys{SystemBase: NewSysBase()})
	return &Client{
		ECS:        ecs,
		Connection: nil,
	}
}

func (Client *Client) Init(){
	Client.ECS.Init()
}

// process all updatable entities
func (Client *Client) Process(){
	Client.ECS.Update()
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
	Client.ECS.Close()
	Client.Connection.Close()
}