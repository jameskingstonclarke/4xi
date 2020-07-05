package src

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type NextTurnEvent struct{
	EventBase
	NewTurn   uint32
}

type State struct {
	*Entity
	*SyncComp
	*StateComp
}

type StateComp struct {
	PlayerID  string
	Turn  	  uint32
	TakenTurn bool
}

// represents the state of the game
// it keeps track of the turn etc
type StateSys struct {
	*SystemBase
	// TODO in theory this should only be a single StateComp
	StateComps []*StateComp
	// the id of the entity that represents us
	OurStateID uint32
}

func (ECS *ECS) CreateState(id string, dirty bool) *State{
	// create a new state with the given id
	state := &State{Entity: ECS.NewEntity("state"), SyncComp: &SyncComp{Dirty: dirty}, StateComp: &StateComp{PlayerID:id}}
	return state
}

// add a new player state to the game
func (ECS *ECS) AddState(state *State) uint32{
	ECS.AddEntity(state.Entity, state.SyncComp, state.StateComp)
	// add the cell to the systems
	for _, system := range ECS.Sys(){
		switch s := system.(type){
		case *NetworkSys:
			s.AddEntity(state.Entity, state.SyncComp)
		case *StateSys:
			s.AddEntity(state.Entity, state.StateComp)
		}
	}
	return state.ID
}

func (S *StateSys) AddEntity(Entity *Entity, StateComp *StateComp){
	S.Entities = append(S.Entities, Entity)
	S.StateComps = append(S.StateComps, StateComp)
	S.Size++
}

func (S *StateSys) Init(){
	S.ECS.RegisterEntity("state", reflect.TypeOf(&State{}), reflect.ValueOf(&State{}).Elem())
	if S.ECS.HostMode == SERVER {
		// at our state to the game
		S.OurStateID = S.ECS.AddState(S.ECS.CreateState("james-host", true))
	}
}

func (S *StateSys) Update(){
	if S.ECS.HostMode & CLIENT != 0 && InputBuffer.KeyPressed == 'n'{
		// player on client side sends a client command to the server indicating it wants the next turn
		S.ECS.Event(ClientCommandEvent{Type: CLIENT_CMD_NEXT_TURN, Side: CLIENT, Data: fmt.Sprintf("{\"id\":%d}", S.OurStateID)})
	}
}
func (S *StateSys) Remove(){}

// testing to see if next turns work
func (S *StateSys) ListenServerCommandEvent(event ServerCommandEvent){
	if event.Side == SERVER && event.Type == SERVER_CMD_NEXT_TURN{

	}
	if event.Side == CLIENT && event.Type == SERVER_CMD_NEXT_TURN{
		CLog("received next turn from server")
		S.ECS.Event(NewWinEvent{
			ID:    "next_turn",
			Title: "next_turn",
			Text: map[string]func(){
				"next turn!":nil,
			},
		})
	}
}

// TODO SERVER
// server listens for commands
func (S *StateSys) ListenClientCommandEvent(event ClientCommandEvent){
	if event.Side == SERVER {
		switch event.Type {
		case CLIENT_CMD_NEXT_TURN:

			SLog("client sent a next turn")
			// TODO first check if the client has taken their turn already as this is spamable
			var result map[string]interface{}
			json.Unmarshal([]byte(event.Data), &result)
			// the client that is taking their turn
			clientID := uint32(result["id"].(float64))
			// used to check if everyone has taken their turn
			turnBuffer := 0
			// check if the client has taken their turn already
			// if they haven't, then take their turn
			for i:=0;i<S.Size;i++{
				if S.Entities[i].ID == clientID{
					if S.StateComps[i].TakenTurn != true{
						S.StateComps[i].TakenTurn = true
						S.StateComps[i].Turn++
					}
				}
				// increase the turn buffer if the client has taken their turn
				if S.StateComps[i].TakenTurn{
					turnBuffer++
				}
			}
			SLog(turnBuffer, S.Size)
			// if everyone has taken their turn, then we go to the next turn
			if turnBuffer == S.Size {
				// update everyone's turns
				for _, client := range S.StateComps{
					client.Turn++
				}

				SLog("sending next turn to client")
				S.ECS.Event(ServerCommandEvent{Side: SERVER, Type: SERVER_CMD_NEXT_TURN})
				S.ECS.Event(ServerCommandEvent{Side: SERVER, Type: SERVER_CMD_SYNC})
			}
		}
	}
}