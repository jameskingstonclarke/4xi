package src


// TODO technically we dont need a player system as we can just use the empire system instead

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

func (P *PlayerStatsComp) Test(){}

func (ECS *ECS) AddPlayer(name string) uint32{
	player := &Player{
		Entity: ECS.NewEntity("player"),
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
	return player.ID
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
}

func (P *PlayerSys) Remove(){
}

func (P *PlayerSys) Close(){

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