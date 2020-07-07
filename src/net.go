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

	// used to broadcast to all clients
	BROADCAST = 0x0
	// directed to one client
	DIRECT    = 0x1
	// used when the server tells the clients it is the next turn
	SERVER_CMD_NEXT_TURN   = 0x98
	// used by the server to indicate clients should sync
	SERVER_CMD_SYNC        = 0x99
	// used by the server to tell a client they are being initialised
	SERVER_CMD_CLIENT_INIT = 0xA0
)

type SyncEntityContainer struct {
	Id float64
	Tag string
	Components []SyncComponentContainer
}

type SyncComponentContainer struct {
	Id string
	Data string
}


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
	ClientConnections map[int]net.Conn
	NumConnections    int
	// address of the server to connect to
	ServerAddress	  string
}

// a sync component is attached to any entity that needs to be synchronized on the network
type SyncComp struct {
	// this is true if the entity is new or has been modified
	Dirty bool
	// TODO for now we don't need this as the component can simply choose not to deserialize itself
	// names of the components that we DON't want to sync with the server, these will be ignored when syncing.
	// this is particularly useful for ignoring client side components e.g. RenderComps etc
	Hidden map[string]struct{}
}

// used by clients to send commands to the server
type ServerCommandEvent struct{
	EventBase
	// whether the command is on the server or client side
	Side uint8
	Type uint32
	Data []byte
	TargetMode uint8
	Target     int
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
	// TODO this is 0 for now, the reason being, the state needs to get the latest state ID
	// after the server has synced
	return 0
	//return 1
}


// TODO SERVER
// listen for client connections. this runs throughout the lifetime of the application
func (N *NetworkSys) PollClientConnections(){
	for Running{
		client, err := N.Listener.Accept()
		if err != nil{
			SLogErr(err)
		}
		clientID := N.NumConnections
		N.ClientConnections[clientID] = client
		N.NumConnections++
		SLog("client connected from ", client.RemoteAddr())
		// as a new client has connected, we need to add a state for the client in the ECS
		N.ECS.AddState(N.ECS.CreateState("client_"+client.RemoteAddr().String(), true))
		// we now need to set the entire ECS system to dirty so the new client
		// can sync the entire server ECS
		for _, sync := range N.SyncComps{
			sync.Dirty = true
		}
		// now sync the new player with every client
		N.ECS.Event(ServerCommandEvent{Side: SERVER, Type: SERVER_CMD_SYNC, TargetMode: BROADCAST})
		// directly tell the client they are being initialised
		N.ECS.Event(ServerCommandEvent{
			Side: SERVER,
			Type: SERVER_CMD_CLIENT_INIT,
			TargetMode: DIRECT,
			Target: clientID,
		})
		SLog("synced new player!")
	}
}


func (N* NetworkSys) Broadcast(command ServerCommandEvent){
	// send a next turn command
	for id, _ := range N.ClientConnections {
		N.Direct(id, command)
	}
}

func (N* NetworkSys) Direct(clientID int, command ServerCommandEvent){
	client, ok := N.ClientConnections[clientID]
	if !ok{
		SLog("cannot find client: ", clientID)
	}
	encoder := json.NewEncoder(client)
	command.Side = CLIENT
	err := encoder.Encode(command)
	if err != nil {
		SLog(err)
	}
}

func (N *NetworkSys) Dispatch(command ServerCommandEvent){
	if command.TargetMode == BROADCAST {
		N.Broadcast(command)
	}else if command.TargetMode == DIRECT{
		N.Direct(command.Target, command)
	}
}
// TODO CLIENT & SERVER
// TODO make this code better jesus
// listen for when the server has done processing and is ready to sync
func (N *NetworkSys) ListenServerCommandEvent(command ServerCommandEvent){
	if command.Side == SERVER && command.Type == SERVER_CMD_CLIENT_INIT{
		N.Dispatch(command)
	}
	if command.Side == SERVER && command.Type == SERVER_CMD_NEXT_TURN{
		N.Dispatch(command)
	}
	if command.Side == SERVER && command.Type == SERVER_CMD_SYNC {
		SLog("server attempting to sync...")
		// TODO
		// we need to iterate over each entity and check if it is dirty. if so, we then get all of it's
		// component data and pack it into the sync command data field. we then serialise this command
		// and broadcast it to the clients.
		entities:="["
		for i := 0; i < N.Size; i++ {
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
		//SLog(syncEvent)
		// we want to send a sync command over the network to all clients
		newCommand := ServerCommandEvent{
			Type: SERVER_CMD_SYNC,
			Side: CLIENT,
			Data: []byte(entities), // TODO this data should be the array of the dirty entities
		}
		N.Dispatch(newCommand)
		SLog("sent sync command to clients ")
	}
	// check if we have received a sync from the server
	if command.Side == CLIENT && command.Type == SERVER_CMD_SYNC {
		CLog("server sent us a sync! ")
		var entities []SyncEntityContainer
		json.Unmarshal(command.Data, &entities)
		// loop through every entity that needs syncing
		for _, e := range entities{
			_, found := N.ECS.Entities[uint32(e.Id)]
			if found{
				// we now need to go through and update the components
				oldComponents := N.ECS.GetEntityComponents(uint32(e.Id))
				// go through each old component and check if it needs updating
				for _, oldComp := range oldComponents{
					// check if it matches with a new comp
					for _, newComp := range e.Components{
						if reflect.TypeOf(oldComp).String()[5:] == newComp.Id{
							var newCompInstance = reflect.New(reflect.TypeOf(oldComp).Elem()).Interface().(Component)
							json.Unmarshal([]byte(newComp.Data), &newCompInstance)
							switch newComp.Id {
							case "PosComp":
								*oldComp.(*PosComp) = *newCompInstance.(*PosComp)
							case "MovementComp":
								*oldComp.(*MovementComp) = *newCompInstance.(*MovementComp)
							case "StateComp":
								*oldComp.(*StateComp) = *newCompInstance.(*StateComp)
							}
						}
					}
				}
			}
			// if our ECS does not contain the entity, then we need to add it to our ECS.
			// this is used for when the server created an entity and needs to sync it with
			// ours
			if !found{
				// create a new instance of the entity
				entityInstance := reflect.New(N.ECS.EntityNameTypeRegistry[e.Tag].Elem()).Interface()
				// now go through and create the actual components for the entity
				var newComponents []Component
				// loop through every old component
				for _, newComponent := range e.Components{
					// check if the name matches the components of the entity
					for _, compType := range N.ECS.EntityComponentTypeRegistry[e.Tag]{
						// check if the components match
						if newComponent.Id == compType.String()[5:] {
							var newCompInstance = reflect.New(compType.Elem()).Interface().(Component)
							json.Unmarshal([]byte(newComponent.Data), &newCompInstance)
							newComponents = append(newComponents, newCompInstance)
						}
					}
				}

				// once we have the components for the entity, we need to set the correct entity fields
				entityVal := reflect.ValueOf(entityInstance).Elem()
				// loop over each field of the new entity
				for i:=0;i<entityVal.NumField();i++{
					// we are setting the actual entity pointer
					if i == 0{
						entityVal.FieldByIndex([]int{0}).Set(reflect.ValueOf(&Entity{
							ID:             uint32(e.Id),
							Tag:            e.Tag,
							ComponentCount: uint32(len(e.Components)),
						}))
					}else{
						// first check if the component is not nil
						// this would be the case for the RenderComp, as the server doesn't generate them,
						// it is up to the client
						component := newComponents[i-1]
						if component != nil{
							entityVal.FieldByIndex([]int{i}).Set(reflect.ValueOf(component))
						}
					}
				}
				// finally add the entity to the correct system
				switch e:=entityInstance.(type){
				case *State:
					N.ECS.AddState(e)
				case *Settlement:
					N.ECS.AddSettlement(e)
				case *Cell:
					N.ECS.AddCell(e)
				case *Unit:
					N.ECS.AddUnit(e)
				}
			}
		}
	}
}

// TODO CLIENT
// listen for systems to send commands to the server
// receive the command, and send it over the connection to the server
func (N *NetworkSys) ListenClientCommandEvent(command ClientCommandEvent){
	if command.Side == CLIENT {
		encoder := json.NewEncoder(N.ServerConnection)
		// this already appends a '\n' character, so we don't need to manually add it
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
	decoder := json.NewDecoder(N.ServerConnection)
	for {
		// create a decoder to listen for server commands
		var command ServerCommandEvent
		// decode the response
		err := decoder.Decode(&command)
		if err != nil{
			CLog(err)
		}
		CLog("received command from the server ", command.Type)
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
		if nextClient != N.NumConnections {
			client, ok := N.ClientConnections[nextClient]
			if !ok{
				SLog("cannot find client with ID ", nextClient)
			}
			nextClient++
			go func() {
				// create a connection and a decoder for each client
				decoder := json.NewDecoder(client)
				for Running {
					// create a decoder to receive the command
					var command ClientCommandEvent
					// decode the response
					err := decoder.Decode(&command)
					if err != nil{
						CLog(err)
					}
					// indicate that the command is now server side
					command.Side = SERVER
					N.ECS.Event(command)
				}
			}()
		}
	}
}