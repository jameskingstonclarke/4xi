package src

import (
	"github.com/nsf/termbox-go"
	"sync"
)

var (
	Running = true
	WaitGroup = sync.WaitGroup{}
)

func Run(){

	ScreenInstance = &Screen{}

	ScreenInstance.Init()
	ScreenInstance.Put('a', 5,1, termbox.AttrBold | termbox.ColorGreen, termbox.ColorWhite)
	ScreenInstance.Put('a', 6,1, termbox.AttrBold | termbox.ColorGreen, termbox.ColorWhite)
	ScreenInstance.Put('a', 7,1, termbox.AttrBold | termbox.ColorGreen, termbox.ColorWhite)

	go ScreenInstance.Draw()
	go ScreenInstance.Poll()


	LogicInstance = &Logic{}
	LogicInstance.Init()

	go LogicInstance.Process()

	WaitGroup.Add(3)
	WaitGroup.Wait()

	ScreenInstance.Close()
	LogicInstance.Close()
}