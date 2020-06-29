package src

import "reflect"

var (
	EId uint32 = 0
)

type ECS struct {
	Systems [][]System
}

func NewECS() *ECS{
	return &ECS{}
}

func (ECS *ECS) Init(){
	for _, s := range ECS.Sys(){
		// check if we can call Initialise on the system
		init, ok := s.(Initialiser)
		if ok{
			init.Init(ECS)
		}
	}
}

func (ECS *ECS) Close(){
	for _, s := range ECS.Sys(){
		// check if we can call Close on the system
		closer, ok := s.(Closer)
		if ok{
			closer.Close(ECS)
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
		s.Update(ECS)
	}
}

func (ECS *ECS) RegisterSystem(System System){
	// check if we can call initialise on the system
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

		//reflect.ValueOf(s).MethodByName("Listen").Call([]reflect.Value{})
		method := reflect.ValueOf(s).MethodByName("Listen")
		valid := method.IsValid()
		if valid{
			// check if the event matches the function type
			if method.Type().In(0) == reflect.TypeOf(Event){
				method.Call([]reflect.Value{reflect.ValueOf(Event)})
			}
		}
	}
}

type Entity struct {
	// Unique id for this entity
	ID   uint32
}

func NewEntity() *Entity{
	e := &Entity{ID: EId}
	EId++
	return e
}

type System interface {
	Update(ECS *ECS)
	Remove()
}

type SystemBase struct {
	Entities []*Entity
	Size     int
	Priority int
}

func NewSysBase() *SystemBase{
	return &SystemBase{
		Entities: nil,
		Size:     0,
		Priority: 0,
	}
}

type Prioritiser interface {
	Priority() int
}

type Initialiser interface {
	Init(ECS *ECS)
}

type Closer interface {
	Close(ECS *ECS)
}

type Event interface {
}

type EventBase struct {
	Entity *Entity
}