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
	PlayerStatsComps []*PlayerStatsComp
	// the current turn we are on
	Turn uint32
	TurnBuffer uint32
}

// player entity
type Player struct{
	*Entity
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
		PlayerStatsComp:   &PlayerStatsComp{
			Name: name,
		},
	}
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *PlayerSys:
			s.AddEntity(player.Entity, player.PlayerStatsComp)
		}
	}
}

func (P *PlayerSys) Init(){
	P.ECS.AddPlayer("james")
}


func (P *PlayerSys) AddEntity(Entity *Entity, PlayerStatsComp *PlayerStatsComp){
	P.Entities = append(P.Entities, Entity)
	P.PlayerStatsComps = append(P.PlayerStatsComps, PlayerStatsComp)
	P.Size++
}

func (P *PlayerSys) Update(){
	// pretend we are the client sending a next turn command
	P.ECS.Event(ClientCommandEvent{})
}

func (P *PlayerSys) Remove(){
}

func (P *PlayerSys) Close(){

}

// TODO SERVER SIDE
// server listens for commands
func (P *PlayerSys) ListenServerCommandEvent(event ServerCommandEvent){
	switch event.Type{
		case CMD_NEXT_TURN:
			P.TurnBuffer++
			// check if everyone has taken their turn
			if P.TurnBuffer == P.Size{
				P.Turn++
				// once everyone has taken their turn, dispatch a next turn event to tell the networksys
				// to dispatch a sync event to all the clients
				P.ECS.Event(ServerProcessEvent{EventBase{}})
			}
	}
}

// TODO CLIENT SIDE
// listen for sync event to update our state
func (P *PlayerSys) ListenSyncEvent(event SyncEvent){
	P.Turn = event.Turn
}