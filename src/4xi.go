package src

import (
	"sync"
)

const (
	CLIENT = 0x0
	HOST   = 0x1
)

var (
	Running = true
	WaitGroup = sync.WaitGroup{}
	LoopMutex = &sync.Mutex{}

	Mode = HOST
)

func Run(){

	ScreenInstance.Init()
	LogicInstance.Init()

	// start the goroutine for the screen events
	go ScreenInstance.Poll()
	go LogicInstance.Process()

	// wait on the screen polling
	WaitGroup.Add(2)
	WaitGroup.Wait()

	ScreenInstance.Close()
	LogicInstance.Close()
}