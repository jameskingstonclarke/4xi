package src

import (
	"bytes"
	"encoding/json"
	"reflect"
)

type ECS struct {
	EId uint32
	Systems [][]System
	// store a reference to each entity
	Entities []*Entity
	// store a reference to the components attached to the entity (linearly added)
	Components []Component
	// used as an index lookup to find the position of the first component of the particular entity
	EntityComponentLookup map[uint32]uint32
	// whether the ECS is a server or client
	HostMode uint8
}

func NewECS(mode uint8) *ECS{
	return &ECS{HostMode: mode, EntityComponentLookup: make(map[uint32]uint32)}
}

// add an entity to the system. this is useful for querying components
// this is likely only used to know which components an entity has
func (ECS *ECS) AddEntity(Entity *Entity, components... Component){
	ECS.Entities = append(ECS.Entities, Entity)
	// set a marker for the component lookup index
	ECS.EntityComponentLookup[Entity.ID] = uint32(len(ECS.Components))
	// TODO this may be wrong, as we might not actually add all the components in this function???
	Entity.ComponentCount = uint32(len(components))
	for _, component := range components{
		ECS.Components = append(ECS.Components, component)
	}
}

// serialize an entity by serializing each component
// TODO make this better, as its gonna be hard to unmarshal the bytes...
func (ECS *ECS) SerializeEntity(id uint32) []byte{
	components:=ECS.GetEntityComponents(id)
	buf := new(bytes.Buffer)
	for _, comp := range components{
		bytes, err := json.Marshal(comp)
		if err != nil{
			SLogErr(err)
		}
		buf.Write(bytes)
	}
	return buf.Bytes()
}

// get all the components attached to a particular entity
// this is likely only used for serialization and de-serialization as the only accessable data this will return
// is the ability to call Serialize() and Deserialize()
func (ECS *ECS) GetEntityComponents(id uint32) []Component{
	// we can just use the id as an index as the id's are incremental
	entity := ECS.Entities[id]
	// get the marker for the first component
	marker := ECS.EntityComponentLookup[id]
	var components []Component
	for i:=0;i<int(entity.ComponentCount);i++{
		// use the marker as the offset, and then we just increase by 1 for each component
		components = append(components, ECS.Components[int(marker)+i])
	}
	return components
}

func (ECS *ECS) Init(){
	for _, s := range ECS.Sys(){
		// check if we can call Initialise on the system
		init, ok := s.(Initialiser)
		if ok{
			init.Init()
		}
	}
}

func (ECS *ECS) Close(){
	for _, s := range ECS.Sys(){
		// check if we can call Close on the system
		closer, ok := s.(Closer)
		if ok{
			closer.Close()
		}
	}
}

func (ECS *ECS) Sys() []System{
	var systems []System
	for _, priority := range ECS.Systems{
		for _, s := range priority{
			systems = append(systems, s)
		}
	}
	return systems
}

func (ECS *ECS) Update(){
	for _, s := range ECS.Sys(){
		s.Update()
	}
}

func (ECS *ECS) RegisterSystem(System System){
	// check if we can prioritise the system
	priority:=0
	var i interface{}=System
	pri, ok:=i.(Prioritiser)
	if ok{
		priority = pri.Priority()
	}
	// check if we need to add more priority layers
	for priority>=len(ECS.Systems){
		ECS.Systems = append(ECS.Systems,nil)
	}
	// add the system to the ECS
	ECS.Systems[priority] = append(ECS.Systems[priority], System)
}

// fire an event into the ECS
func (ECS *ECS) Event(Event Event){
	// check each system to see if it is capable of hearing the event
	for _, s := range ECS.Sys(){
		method := reflect.ValueOf(s).MethodByName("Listen"+reflect.TypeOf(Event).String()[4:])
		valid := method.IsValid()
		if valid{
			method.Call([]reflect.Value{reflect.ValueOf(Event)})
		}
	}
}

type Entity struct {
	// Unique id for this entity
	ID   		   uint32
	// used as a util for the ECS
	ComponentCount uint32
}

func (ECS *ECS) NewEntity() *Entity{
	e := &Entity{ID: ECS.EId}
	ECS.EId++
	return e
}

type System interface {
	// update called every frame
	Update()
	Remove()
}

type SystemBase struct {
	ECS		 *ECS
	Entities []*Entity
	Size     int
	Priority int
}

func NewSysBase(ECS *ECS) *SystemBase{
	return &SystemBase{
		ECS: 	  ECS,
		Entities: nil,
		Size:     0,
		Priority: 0,
	}
}

type Prioritiser interface {
	Priority() int
}

type Initialiser interface {
	Init()
}

type Closer interface {
	Close()
}

type Component interface {
	//Serialize() []byte
}

type Event interface {
}

type EventBase struct {
}