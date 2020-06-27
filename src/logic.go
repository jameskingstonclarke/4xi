package src

import "github.com/nsf/termbox-go"

type Logic struct {
}

var (
	LogicInstance = &Logic{}
)

func (Logic *Logic) Process(){
	for Running{
		ScreenInstance.Text("hello world", 5,1, termbox.AttrBold | termbox.ColorGreen, termbox.ColorWhite)
	}
}

func (Logic *Logic) Init(){
}

func (Logic *Logic) Close(){
	WaitGroup.Done()
}