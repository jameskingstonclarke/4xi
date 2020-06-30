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
type ServerCommandEvent struct{
	*EventBase
	Type uint32
}

// used by clients to send commands to the server
type ClientCommandEvent struct{
	*EventBase
	Type uint32
}

// used by the server to sync data to clients
type SyncEvent struct {
	EventBase
	Turn uint32
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
		// start listening for syncs from the server
		go N.ListenForSync()
	}else if Mode & SERVER != 0 {
		// start listening for client commands
		go N.ListenForClientCommands()
	}
}
func (N *NetworkSys) Remove(){}
func (N *NetworkSys) Close(){}


// TODO SERVER
// listen for when the server has done processing and is ready to sync
func (N *NetworkSys) ListenServerProcessEvent(event ServerProcessEvent){
	// send out a sync to all clients
	for _, client := range N.ClientConnections{
		// send the sync event over the network
		// create a buffer for the event
		buf := new(bytes.Buffer)
		// create an encoder object
		gobobj := gob.NewEncoder(buf)
		// write the event to the buffer (new sync event)
		gobobj.Encode(SyncEvent{EventBase: EventBase{}, Turn: 0})
		// send the event to the server
		client.Write(buf.Bytes())
	}
}

// TODO CLIENT
// listen for systems to send commands to the server
// recieve the command, and send it over the connection to the server
func (N *NetworkSys) ListenClientCommandEvent(event ClientCommandEvent){
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

// TODO CLIENT
// listen for sync messages from the server connection, once we recieve a sync over the network
// dispatch the sync locally
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

// TODO server listens for commands from clients
func (N *NetworkSys) ListenForClientCommands(){
	for _, client := range N.ClientConnections{
		go func() {
			tmp := make([]byte, 500)
			for {
				client.Read(tmp)
				// convert bytes into Buffer
				tmpbuff := bytes.NewBuffer(tmp)
				event := new(ServerCommandEvent)
				// creates a decoder object
				gobobj := gob.NewDecoder(tmpbuff)
				// decodes the buffer into the CommandEvent struct
				gobobj.Decode(event)

				// now we have the command event, we need to update our server game state
				// how do we do this... do we broadcast a server-side sync event??? no clue m8

				// TODO for now we broadcast the command server side for systems to recieve
				N.ECS.Event(event)
			}
		}()
	}
}