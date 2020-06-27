package src

import "github.com/nsf/termbox-go"

type Logic struct {
}

var (
	LogicInstance = &Logic{}
)

func (Logic *Logic) Process(){
	for Running{
		// lock the screen, and then process all logic
		ScreenMutex.Lock()
		ScreenInstance.Text("hello world", 5,1, termbox.AttrBold | termbox.ColorGreen, termbox.ColorWhite)
		ScreenInstance.Rect('#', 1,1, 10,10, termbox.AttrBold | termbox.ColorGreen, 0, true)
		// finally release the screen for the game to render
		ScreenMutex.Unlock()
	}
}

func (Logic *Logic) Init(){
}

func (Logic *Logic) Close(){
	WaitGroup.Done()
}