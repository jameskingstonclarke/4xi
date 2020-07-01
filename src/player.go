package src

// used by the server to tell the networks to dispatch a sync event as we have
// moved onto the next turn
// (syncs happen every turn)
type ServerProcessEvent struct {
	EventBase
}

// player system
type PlayerSys struct {
	*SystemBase
	SyncComps		 []*SyncComp
	PlayerStatsComps []*PlayerStatsComp
	// the current turn we are on
	Turn 	   int
	TurnBuffer int
	Done bool
}

// player entity
type Player struct{
	*Entity
	*SyncComp
	*PlayerStatsComp
}

// component for storing player info
type PlayerStatsComp struct {
	// name of the player
	Name  string
}

func (ECS *ECS) AddPlayer(name string){
	player := &Player{
		Entity: NewEntity(),
		SyncComp: &SyncComp{Dirty: false},
		PlayerStatsComp:   &PlayerStatsComp{
			Name: name,
		},
	}

	ECS.AddEntity(player.Entity, player.SyncComp, player.PlayerStatsComp)

	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *NetworkSys:
			s.AddEntity(player.Entity, player.SyncComp)
		case *PlayerSys:
			s.AddEntity(player.Entity, player.PlayerStatsComp, player.SyncComp)
		}
	}
}

func (P *PlayerSys) Init(){
	P.ECS.AddPlayer("james")
}

func (P *PlayerSys) AddEntity(Entity *Entity, PlayerStatsComp *PlayerStatsComp, SyncComp *SyncComp){
	P.Entities = append(P.Entities, Entity)
	P.PlayerStatsComps = append(P.PlayerStatsComps, PlayerStatsComp)
	P.SyncComps = append(P.SyncComps, SyncComp)
	P.Size++
}

func (P *PlayerSys) Update(){
	if P.ECS.HostMode & CLIENT != 0 && !P.Done && InputBuffer.KeyPressed == 'n'{
		// player on client side sends a client command to the server indicating it wants the next turn
		P.ECS.Event(ClientCommandEvent{Type: CLIENT_CMD_NEXT_TURN, Side: CLIENT})
		P.Done = true
	}

	// used for testing so we can set the dirty flag on the SyncComp
	if P.ECS.HostMode & CLIENT != 0 && InputBuffer.KeyPressed == 'p'{
		P.SyncComps[0].Dirty = true
	}
}

func (P *PlayerSys) Remove(){
}

func (P *PlayerSys) Close(){

}

// TODO SERVER SIDE
// server listens for commands
func (P *PlayerSys) ListenClientCommandEvent(event ClientCommandEvent){
	if event.Side == SERVER {
		switch event.Type {
		case CLIENT_CMD_NEXT_TURN:
			P.TurnBuffer++
			// check if everyone has taken their turn
			if P.TurnBuffer == P.Size {
				P.Turn++
				// once everyone has taken their turn, dispatch a next turn event to tell the networksys
				// to dispatch a sync event to all the clients
				// TODO this is the problem for some reason, obviously syncing is fucked
				SLog("client next turn'd")
				P.ECS.Event(ServerCommandEvent{Type: SERVER_CMD_SYNC, Side: SERVER})
			}
		}
	}
}

// TODO CLIENT
// listen for sync event to update our state
func (P *PlayerSys) ListenServerCommandEvent(event ServerCommandEvent){
	if event.Side == CLIENT{
		switch event.Type{
		case SERVER_CMD_SYNC:
			CLog("server sent us a sync! ")
		}
	}
}