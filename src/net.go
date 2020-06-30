package src

import (
	"bytes"
	"encoding/gob"
	"net"
)

const (
	CLIENT = 0x1 << 0
	SERVER = 0x1 << 1

	CMD_NEXT_TURN = 0x0
)

type NetworkSys struct {
	*SystemBase
	NetworkComps []*NetworkComp
	// whether the system is a server or client
	Mode uint8

	// connection to the server (used if we are a client)
	ServerConnection net.Conn
	// connections to clients (used if we are a server)
	ClientConnections []net.Conn
}

type NetworkComp struct {

}

// used by clients to send commands to the server
type CommandEvent struct{
	*EventBase
	Type uint32
}

// used by the server to sync data to clients
type SyncEvent struct {
	*EventBase
}

func (N *NetworkSys) Init(){
	if N.Mode & CLIENT != 0{
		// if we are the client, start a goroutine for listening for server syncs
	}else if N.Mode & SERVER != 0{
		// if we are the server, start a goroutine for listening for client commands
	}
}
func (N *NetworkSys) Update(){
	// if we are the client, we listen for messages for the server connection,
	// when we retrieve these messages, we broadcast a sync event with the data the
	// server sent us.
	// if we are the server listen for client commands via client connections. Once we
	// have received all commands, send a message to the clients to trigger a sync event.
	if Mode & CLIENT != 0{
		go N.ListenForSync()
	}else if Mode & SERVER != 0 {
		go N.ListenForCommands()
	}
}
func (N *NetworkSys) Remove(){}
func (N *NetworkSys) Close(){}


// TODO Note, the netsys does NOT listen for sync events. this is the job of individual systems to do
// listen for systems that are sending command events.
// when we recieve these command events, we send them to the server command queue directly.
func (N *NetworkSys) Listen(event CommandEvent){
	// send the command event over the network
	// create a buffer for the event
	buf := new(bytes.Buffer)
	// create an encoder object
	gobobj := gob.NewEncoder(buf)
	// write the event to the buffer
	gobobj.Encode(event)
	// send the event to the server
	N.ServerConnection.Write(buf.Bytes())
}



// listen for sync messages from the server connection
func (N *NetworkSys) ListenForSync(){
	// create a temp buffer for the sync events
	tmp := make([]byte, 500)
	for {
		N.ServerConnection.Read(tmp)
		// convert bytes into Buffer
		tmpbuff := bytes.NewBuffer(tmp)
		event := new(SyncEvent)
		// creates a decoder object
		gobobj := gob.NewDecoder(tmpbuff)
		// decodes the buffer into the SynEvent struct
		gobobj.Decode(event)
		// broadcast the sync event
		N.ECS.Event(event)
	}
}

// listen for commands from clients
func (N *NetworkSys) ListenForCommands(){

}