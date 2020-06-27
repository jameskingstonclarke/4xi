package src

import (
	"github.com/nsf/termbox-go"
)

const (
	DEFAULT_WIDTH = 100
	DEFAULT_HEIGHT = 50
)

type Screen struct {
	Width, Height int
	CellBuffer [][]termbox.Cell
}

var (
	ScreenInstance *Screen
)

func (Screen *Screen) Init(){
	err := termbox.Init()
	if err != nil{
		LogErr(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	Screen.Width = DEFAULT_WIDTH
	Screen.Height = DEFAULT_HEIGHT
	Screen.Resize()
}

func (Screen *Screen) Resize(){
	Screen.CellBuffer = make([][]termbox.Cell, Screen.Width)
	for i := range Screen.CellBuffer {
		Screen.CellBuffer[i] = make([]termbox.Cell, Screen.Height)
	}
}

func (Screen *Screen) Put(r rune, x, y int, bg, fg termbox.Attribute){
	Screen.CellBuffer[x][y].Ch = r
	Screen.CellBuffer[x][y].Fg = fg
	Screen.CellBuffer[x][y].Bg = bg
}

func (Screen *Screen) Draw(){
	for Running {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for y := 0; y < Screen.Height; y++ {
			for x := 0; x < Screen.Width; x++ {
				cell := Screen.CellBuffer[x][y]
				termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
			}
		}
		termbox.Flush()
	}
}

func (Screen *Screen) Poll() {
	defer WaitGroup.Done()
	for Running {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				Running = false
			case termbox.KeyF1:
				Screen.Put('k', ev.MouseX, ev.MouseY, termbox.ColorGreen, termbox.AttrBold)
			default:
				break
			}
		case termbox.EventMouse:
			Screen.Put('m', ev.MouseX, ev.MouseY, 0, 0)
			break
		case termbox.EventResize:
			Screen.Width = ev.Width
			Screen.Height = ev.Height
			Screen.Resize()
		}
	}
}

func (Screen *Screen) Close(){
	termbox.Close()
	WaitGroup.Done()
}