package src

import (
	"sync"
)

var (
	Running = true
	WaitGroup = sync.WaitGroup{}
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