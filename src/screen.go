package src

import (
	"github.com/nsf/termbox-go"
	"sync"
)

const (
	DEFAULT_WIDTH = 200
	DEFAULT_HEIGHT = 100
)

type Screen struct {
	Width, Height int
	CellBuffer [][]termbox.Cell

}

var (
	ScreenInstance = &Screen{}
	ScreenMutex = &sync.Mutex{}
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

func (Screen *Screen) Text(text string, x, y int, fg, bg termbox.Attribute){
	for i, r := range text {
		if x+i >= Screen.Width {
			x=0
			y++
		}
		Screen.CellBuffer[x+i][y].Ch = r
		Screen.CellBuffer[x+i][y].Fg = fg
		Screen.CellBuffer[x+i][y].Bg = bg
	}
}

func (Screen *Screen) Char(r rune, x, y int, fg, bg termbox.Attribute){
	Screen.CellBuffer[x][y].Ch = r
	Screen.CellBuffer[x][y].Fg = fg
	Screen.CellBuffer[x][y].Bg = bg
}

func (Screen *Screen) Rect(r rune, x, y, width, height int, fg, bg termbox.Attribute, fill bool){
	for xTmp:=x; xTmp<x+width; xTmp++ {
		for yTmp:=y; yTmp<y+height; yTmp++ {
			if fill {
				Screen.CellBuffer[xTmp][yTmp].Ch = r
				Screen.CellBuffer[xTmp][yTmp].Fg = fg
				Screen.CellBuffer[xTmp][yTmp].Bg = bg
			}else{
				if (xTmp==x || xTmp==x+width-1) || (yTmp==y || yTmp==y+height-1) {
					Screen.CellBuffer[xTmp][yTmp].Ch = r
					Screen.CellBuffer[xTmp][yTmp].Fg = fg
					Screen.CellBuffer[xTmp][yTmp].Bg = bg
				}
			}
		}
	}
}

func (Screen *Screen) Draw(){
	for Running {
		ScreenMutex.Lock()
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for y := 0; y < Screen.Height; y++ {
			for x := 0; x < Screen.Width; x++ {
				cell := Screen.CellBuffer[x][y]
				termbox.SetCell(x, y, cell.Ch, cell.Fg, cell.Bg)
			}
		}
		termbox.Flush()
		ScreenMutex.Unlock()
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
				Screen.Char('k', ev.MouseX, ev.MouseY, termbox.ColorGreen, termbox.AttrBold)
			default:
				break
			}
		case termbox.EventMouse:
			Screen.Char('m', ev.MouseX, ev.MouseY, 0, 0)
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