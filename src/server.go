package src

import (
	"bufio"
	"net"
)

// represents a game server
type Server struct {
	Listener  net.Listener
	Clients   map[uint8]net.Conn
	Players   []*Player
	GameState *GameState
}

func (Server *Server) Process(){

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

	for _, entity := range Server.GameState.Entities{
		entity.Update()
	}
}

func (Server *Server) Init(){
	listener, err := net.Listen("tcp", "localhost:7777")
	if err != nil{
		LogErr(err)
	}
	Server.Listener = listener
}

func (Server *Server) AcceptClient(){
	conn, err := Server.Listener.Accept()
	if err != nil{
		LogErr(err)
	}
	Server.Clients[uint8(len(Server.Clients)+1)] = conn
}

func (Server *Server) BroadcastMsg(msg string){
	for client, _ := range Server.Clients {
		Server.SendMsg(client, msg)
	}
}

// send a message to a particular client
func (Server *Server) SendMsg(client uint8, msg string){
	Server.Clients[client].Write([]byte(msg+"\n"))
}

// request (listen) for a message from a particular client
func (Server *Server) ListenMsg(client uint8) string{
	msg, err := bufio.NewReader(Server.Clients[client]).ReadString('\n')
	if err != nil{
		LogErr(err)
	}
	return msg
}

func (Server *Server) Close(){
	Server.Listener.Close()
	for _, c := range Server.Clients{
		c.Close()
	}
}