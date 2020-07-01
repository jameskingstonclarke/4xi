package src

import (
	"bytes"
	"encoding/gob"
	"net"
)

const (
	CLIENT = 0x1 << 1
	SERVER = 0x1 << 2
	PORT   = ":7777"
	// used by the client to indicate it wants to go to the next turn
	CLIENT_CMD_NEXT_TURN = 0x0
	// used by the server to indicate clients should sync
	SERVER_CMD_SYNC      = 0x1
)

// this system handles everything network related. it can be in 1 of 2 modes, CLIENT or SERVER. the
// mode determines the behaviour and will send and listen for the relevant events.
// the system stores all the sync components. these components are sent across the network
type NetworkSys struct {
	*SystemBase
	SyncComps []*SyncComp
	// connection to the server (used if we are a client)
	ServerConnection  net.Conn
	// connections to clients (used if we are a server)
	Listener          net.Listener
	ClientConnections []net.Conn
	// address of the server to connect to
	ServerAddress	  string
}

// a sync component is attached to any entity that needs to be synchronized on the network
type SyncComp struct {
	// this is true if the entity's state has been changed. if this is true, the entity is synced across the network
	Dirty bool
}

func (S *SyncComp) Serialize() []byte{
	return []byte{0}
}

// used by clients to send commands to the server
type ServerCommandEvent struct{
	EventBase
	// whether the command is on the server or client side
	Side uint8
	Type uint32
	Data []byte
}

// used by clients to send commands to the server
type ClientCommandEvent struct{
	EventBase	// whether the command is on the server or client side
	Side uint8
	Type uint32
	Data []byte
}

// used by the server to sync data to clients
type SyncEvent struct {
	EventBase
	DirtyEntities []Entity
}

type SyncAction struct {}

func (N *NetworkSys) AddEntity(Entity *Entity, SyncComp *SyncComp){
	N.Entities = append(N.Entities, Entity)
	N.SyncComps = append(N.SyncComps, SyncComp)
	N.Size++
}

func (N *NetworkSys) Init(){
	// if we are the client, we listen for messages for the server connection,
	// when we retrieve these messages, we broadcast a sync event with the data the
	// server sent us.
	// if we are the server listen for client commands via client connections. Once we
	// have received all commands, send a message to the clients to trigger a sync event.
	if N.ECS.HostMode & CLIENT != 0{
		CLog("client attempting to connect ", N.ServerAddress+PORT)
		conn, err := net.Dial("tcp", N.ServerAddress+PORT)
		if err != nil{
			CLogErr(err)
		}
		N.ServerConnection = conn
		CLog("client connected to ", N.ServerAddress+PORT)
		// start listening for commands from the server
		go N.PollServerCommands()
	}else if N.ECS.HostMode & SERVER != 0{
		SLog("server listening on port ", PORT)
		listener, err := net.Listen("tcp", PORT)
		if err != nil{
			SLogErr(err)
		}
		N.Listener = listener


		// TODO we need a way to concurrently poll connections AND

		go N.PollClientConnections()
		// start listening for client commands
		go N.PollClientCommands()
	}
}
func (N *NetworkSys) Update(){}
func (N *NetworkSys) Remove(){}
func (N *NetworkSys) Close(){}

// the network code must be ran AFTER every other system has updated their state.
// this allows the server to queue up the changes, and then send them across the network
func (N *NetworkSys) Priority() int {
	return 1
}


// TODO SERVER
// listen for client connections. this runs throughout the lifetime of the application
func (N *NetworkSys) PollClientConnections(){
	for Running{
		client, err := N.Listener.Accept()
		if err != nil{
			SLogErr(err)
		}
		N.ClientConnections = append(N.ClientConnections, client)
		SLog("client connected from ", client.RemoteAddr())
	}
}

// TODO SERVER
// listen for when the server has done processing and is ready to sync
func (N *NetworkSys) ListenServerCommandEvent(event ServerCommandEvent){
	if event.Side == SERVER && event.Type == SERVER_CMD_SYNC {
		// send out a sync to all clients
		for _, client := range N.ClientConnections {
			// TODO fix this
			//syncEvent := SyncEvent{
			//	EventBase:     EventBase{},
			//	DirtyEntities: nil,
			//}
			//// store all the entities that need synchronizing in the SyncEvent
			//// iterate over each entity that needs to be synced
			//for i := 0; i < N.Size; i++ {
			//	// if the entity is dirty (it has been changed), it needs synchronizing
			//	if N.SyncComps[i].Dirty {
			//		// clear the dirty flag ready for the next check
			//		// if we didn't clear it, every frame the server would attempt to sync the entity
			//		N.SyncComps[i].Dirty = false
			//		// serialize the entity
			//		serial := N.ECS.SerializeEntity(N.Entities[i].ID)
			//		SLog(serial)
			//	}
			//}
			//SLog(syncEvent)
			// we want to send a sync command over the network to all clients
			command := ServerCommandEvent{
				Type: SERVER_CMD_SYNC,
				Side: SERVER,
				Data: nil,
			}
			// send the sync event over the network
			// create a buffer for the event
			buf := new(bytes.Buffer)
			// create an encoder object
			gobobj := gob.NewEncoder(buf)
			// write the event to the buffer (new sync event)
			gobobj.Encode(command)
			// send the event to the server
			client.Write(buf.Bytes())
		}
	}
}

// TODO CLIENT
// listen for systems to send commands to the server
// recieve the command, and send it over the connection to the server
func (N *NetworkSys) ListenClientCommandEvent(event ClientCommandEvent){
	if event.Side == CLIENT {
		// send the command event over the network
		// create a buffer for the event
		buf := new(bytes.Buffer)
		// create an encoder object
		gobobj := gob.NewEncoder(buf)
		// write the event to the buffer
		gobobj.Encode(event)
		// send the event to the server
		_, err := N.ServerConnection.Write(buf.Bytes())
		if err != nil {
			CLogErr(err)
		}
		CLog("sent command ", event, " to server")
	}
}

// TODO CLIENT
// listen for sync messages from the server connection, once we receive a sync over the network
// dispatch the sync locally
func (N *NetworkSys) PollServerCommands(){
	// create a temp buffer for the sync events
	tmp := make([]byte, 500)
	for {
		N.ServerConnection.Read(tmp)
		// convert bytes into Buffer
		tmpbuff := bytes.NewBuffer(tmp)
		// TODO for some reason we are losing the command type
		command := ServerCommandEvent{}
		// creates a decoder object
		gobobj := gob.NewDecoder(tmpbuff)
		// decodes the buffer into the SynEvent struct
		gobobj.Decode(command)
		CLog("received command from the server ", command)
		// indicate that the command is now client side
		command.Side = CLIENT
		// finally broadcast the server command so systems can react
		N.ECS.Event(command)
	}
}

// TODO SERVER
// server listens for commands from clients
func (N *NetworkSys) PollClientCommands(){
	nextClient := 0
	for Running{
		// this allows us to dispatch a new command listener for the next incoming client
		if nextClient != len(N.ClientConnections){
			client := N.ClientConnections[nextClient]
			nextClient++
			go func() {
				tmp := make([]byte, 500)
				for Running {
					client.Read(tmp)
					// convert bytes into Buffer
					tmpbuff := bytes.NewBuffer(tmp)
					command := ClientCommandEvent{}
					// creates a decoder object
					gobobj := gob.NewDecoder(tmpbuff)
					SLog("received command from client ", command)
					// indicate that the command is now server side
					command.Side = SERVER
					// decodes the buffer into the CommandEvent struct
					gobobj.Decode(command)
					N.ECS.Event(command)
				}
			}()
		}
	}
}