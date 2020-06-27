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

	go ScreenInstance.Draw()
	go ScreenInstance.Poll()

	LogicInstance.Init()

	go LogicInstance.Process()

	WaitGroup.Add(3)
	WaitGroup.Wait()

	ScreenInstance.Close()
	LogicInstance.Close()
}