package src

import (
	"encoding/json"
	"net"
	"reflect"
)

const (
	CLIENT = 0x1 << 1
	SERVER = 0x1 << 2
	PORT   = ":7777"

	// used by the client to indicate it wants to go to the next turn
	CLIENT_CMD_NEXT_TURN = 0x1
	// used by the client to tell the server it is moving a unit
	CLIENT_CMD_MOVE_UNIT = 0x2

	// used when the server tells the clients it is the next turn
	SERVER_CMD_NEXT_TURN = 0x98
	// used by the server to indicate clients should sync
	SERVER_CMD_SYNC      = 0x99
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
	// TODO for now we don't need this as the component can simply choose not to deserialize itself
	// names of the components that we DON't want to sync with the server, these will be ignored when syncing.
	// this is particularly useful for ignoring client side components e.g. RenderComps etc
	Hidden []string
}

func (S *SyncComp) Deserialize(data interface{}){

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
	Data string
	//Data []byte
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

// TODO CLIENT & SERVER
// listen for when the server has done processing and is ready to sync
func (N *NetworkSys) ListenServerCommandEvent(command ServerCommandEvent){
	// TODO make this code better jesus
	if command.Side == SERVER && command.Type == SERVER_CMD_NEXT_TURN{
		// send a next turn command
		for _, client := range N.ClientConnections {
			encoder := json.NewEncoder(client)
			command.Side = CLIENT
			err := encoder.Encode(command)
			if err != nil{
				SLog(err)
			}
		}
	}
	if command.Side == SERVER && command.Type == SERVER_CMD_SYNC {
		SLog("server attempting to sync...")
		// TODO
		// we need to iterate over each entity and check if it is dirty. if so, we then get all of it's
		// component data and pack it into the sync command data field. we then serialise this command
		// and broadcast it to the clients.
		entities:="["
		for i := 0; i < N.Size; i++ {
			SLog("checking entity...")
			// if the entity is dirty (it has been changed), it needs synchronizing
			if N.SyncComps[i].Dirty {
				// clear the dirty flag ready for the next check
				// if we didn't clear it, every frame the server would attempt to sync the entity
				N.SyncComps[i].Dirty = false
				// serialize the entity
				serial := N.ECS.SerializeEntity(N.Entities[i].ID)
				entities+=serial
				entities+=","
			}
		}
		entities = entities[:len(entities)-1]
		entities+="]"
		SLog(entities)
		//SLog(syncEvent)
		// we want to send a sync command over the network to all clients
		newCommand := ServerCommandEvent{
			Type: SERVER_CMD_SYNC,
			Side: CLIENT,
			Data: []byte(entities), // TODO this data should be the array of the dirty entities
		}
		// send out a sync to all clients
		for _, client := range N.ClientConnections {
			encoder := json.NewEncoder(client)
			err := encoder.Encode(newCommand)
			if err != nil{
				SLog(err)
			}
		}
		SLog("sent sync command to clients ")
	}
	// check if we have recieved a sync from the server
	if command.Side == CLIENT && command.Type == SERVER_CMD_SYNC {
		CLog("server sent us a sync! ", string(command.Data))
		// first get the map of entities that need syncing.
		// this is an array of maps, a key for the entity ID and the value being the components.
		// NOTE currently every component is sync'ed, even ones that aren't changed and ones that we don't want to.
		// in the new system, we will have a slice of ones to ignore and these will not be sent across the network.
		var result []interface{}
		json.Unmarshal(command.Data, &result)
		// gor through each entity and check if we need to update the entity in this system
		for _, encodedEntity := range result {
			entity := encodedEntity.(map[string]interface{})
			// get the id of the entity, we now need to check over each entity in the ECS and check if it matches
			id := uint32(entity["id"].(float64))
			for i:=0;i<N.ECS.Size;i++{
				// we found a match
				if N.ECS.Entities[i].ID == id{
					// get the components of the entity we need to update
					oldComponents := N.ECS.GetEntityComponents(id) // old components
					// unmarshal the entity that we have been sent
					newComponents := entity["components"].([]interface{}) // this is an array of maps
					for _, c := range newComponents{
						comp := c.(map[string]interface{})
						compID := reflect.ValueOf(comp).MapKeys()[0].String()
						compValue := comp[compID]
						// now we need to marshal the compValue into the correct component, and then
						// update the correct component in the oldComponents slice.
						// if we find the matching component then deserialize it.
						for _, oldComp := range oldComponents{
							if reflect.TypeOf(oldComp).String()[5:] == compID{
								oldComp.Deserialize(compValue)
							}
						}
					}
				}
			}
		}
	}
}

// TODO CLIENT
// listen for systems to send commands to the server
// recieve the command, and send it over the connection to the server
func (N *NetworkSys) ListenClientCommandEvent(command ClientCommandEvent){
	if command.Side == CLIENT {
		encoder := json.NewEncoder(N.ServerConnection)
		err := encoder.Encode(command)
		if err != nil {
			CLogErr(err)
		}
		CLog("sent command ", command, " to server")
	}
}

// TODO CLIENT
// listen for sync messages from the server connection, once we receive a sync over the network
// dispatch the sync locally
func (N *NetworkSys) PollServerCommands(){
	for {
		// create a decoder to listen for server commands
		decoder := json.NewDecoder(N.ServerConnection)
		var command ServerCommandEvent
		// decode the response
		err := decoder.Decode(&command)
		if err != nil{
			CLog(err)
		}
		//CLog("received command from the server ", command)
		// indicate that the command is now client side
		command.Side = CLIENT
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
				for Running {
					// create a decoder to receive the command
					decoder := json.NewDecoder(client)
					var command ClientCommandEvent
					// decode the response
					err := decoder.Decode(&command)
					if err != nil{
						CLog(err)
					}
					//SLog("received command from client ", command)
					// indicate that the command is now server side
					command.Side = SERVER
					N.ECS.Event(command)
				}
			}()
		}
	}
}