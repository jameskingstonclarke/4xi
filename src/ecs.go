package src

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
)

type ECS struct {
	// number of entities
	Size int
	// store a map of each priority level containing systems
	Systems [][]System
	// store a reference to each entity
	Entities map[uint32]*Entity
	// we store a type registry to store the type of a component
	EntityNameTypeRegistry map[string]reflect.Type	// we store a type registry to store the type of a component
	EntityComponentTypeRegistry map[string][]reflect.Type
	// store a reference to the components attached to the entity (linearly added)
	Components []Component
	// used as an index lookup to find the position of the first component of the particular entity
	EntityComponentLookup map[uint32]uint32
	// whether the ECS is a server or client
	HostMode uint8
	EventMut sync.Mutex
}

func NewECS(mode uint8) *ECS{
	return &ECS{
		HostMode: mode,
		Entities: make(map[uint32]*Entity),
		EntityComponentLookup: make(map[uint32]uint32),
		EntityNameTypeRegistry: make(map[string]reflect.Type),
		EntityComponentTypeRegistry: make(map[string][]reflect.Type),
	}
}

// register an entity by adding its type, and components to the ECS
func (ECS *ECS) RegisterEntity(tag string, t reflect.Type, components reflect.Value){
	// register the entity type to the ECS
	ECS.EntityNameTypeRegistry[tag] = t
	// get the types of the components
	var types []reflect.Type
	for i := 1; i < components.NumField(); i++ {
		fieldType := components.Type().Field(i).Type
		types = append(types, fieldType)
	}
	// register the type of the components to the ECS
	ECS.EntityComponentTypeRegistry[tag] = types
}

// add an entity to the system. this is useful for querying components
// this is likely only used to know which components an entity has
func (ECS *ECS) AddEntity(Entity *Entity, components... Component){
	_, exists := ECS.Entities[Entity.ID]
	if exists{
		SLog("entity with ID ", Entity.ID, " already exists!")
	}
	ECS.Entities[Entity.ID] = Entity
	// set a marker for the component lookup index
	ECS.EntityComponentLookup[Entity.ID] = uint32(len(ECS.Components))
	// TODO this may be wrong, as we might not actually add all the components in this function???
	Entity.ComponentCount = uint32(len(components))
	for _, component := range components{
		ECS.Components = append(ECS.Components, component)
	}
	ECS.Size++
}

func (ECS *ECS) GetEntity(id uint32) *Entity{
	for _, e := range ECS.Entities{
		if e.ID == id{
			return e
		}
	}
	return nil
}

// serialize an entity by serializing each component
// TODO make this better, as its gonna be hard to unmarshal the bytes...
func (ECS *ECS) SerializeEntity(id uint32) string{
	entity := ECS.GetEntity(id)
	components:=ECS.GetEntityComponents(id)
	buf := fmt.Sprintf("{\"Id\":%d, \"Tag\":\"%s\", \"Components\":[",id, entity.Tag)
	for _, comp := range components{
		compID := reflect.TypeOf(comp).String()[5:]
		marshal, err := json.Marshal(&comp)
		if err != nil{
			SLogErr(err)
		}
		newMarshal := strings.Replace(string(marshal), "\"", "\\\"", -1)
		comp := fmt.Sprintf("{\"Id\":\"%s\", \"Data\":\"%s\"},", compID, newMarshal)
		buf+=comp
	}
	buf = buf[:len(buf)-1] // remove the last ','
	buf+="]}"
	return buf
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
	// used to tag the entity with a string identifier
	Tag            string
	// used as a util for the ECS
	ComponentCount uint32
}

func (ECS *ECS) NewEntity(tag string) *Entity{
	e := &Entity{ID: uint32(rand.Intn(10000)), Tag: tag}
	_, found:=ECS.Entities[e.ID]
	for found==true {
		e.ID = uint32(rand.Intn(10000))
		_, found=ECS.Entities[e.ID]
	}
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

type Component interface {}

type Event interface {
}

type EventBase struct {
}