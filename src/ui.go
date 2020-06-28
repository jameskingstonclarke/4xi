package src

import (
	"github.com/gdamore/tcell"
)

type UI interface {
	Draw()
	Update()
}

type UITemplate struct {
	Enabled bool
}

func (UITemplate *UITemplate) Enable(enabled bool){
	UITemplate.Enabled = enabled
}

type Callback func()

type Text struct {
	UITemplate
	T   	 string
	Pos 	 Vec
	Callback Callback
	Style    tcell.Style
}

func (Text *Text) Update(){
	if Text.Enabled {
		if ScreenInstance.InputBuffer.MousePressed&tcell.Button1 != 0 {
			if ScreenInstance.InputBuffer.MousePos.X > Text.Pos.X && ScreenInstance.InputBuffer.MousePos.X < Text.Pos.X+len(Text.T) {
				if ScreenInstance.InputBuffer.MousePos.Y == Text.Pos.Y {
					Text.Callback()
				}
			}
		}
	}
}

func (Text *Text) Draw(){
	if Text.Enabled {
		ScreenInstance.Text(Text.T, Text.Pos, tcell.StyleDefault)
	}
}

type Window struct {
	UITemplate
	Title string
	Pos, Size  Vec
}

func (Window *Window) Update(){
	if Window.Enabled{
		if ScreenInstance.InputBuffer.MousePressed & tcell.Button1 != 0{
			if ScreenInstance.InputBuffer.MousePos.Equals(Window.Pos.Add(V2(Window.Size.X,0))){
				Window.Enable(false)
			}
		}
	}
}

func (Window *Window) Draw(){
	if Window.Enabled {
		// draw the main body
		for col := Window.Pos.X; col < Window.Pos.X+Window.Size.X; col++ {
			for row := Window.Pos.Y; row < Window.Pos.Y+Window.Size.Y; row++ {
				ScreenInstance.Char(tcell.RuneBlock, V2(col, row), tcell.StyleDefault)
			}
		}
		// draw the left & rightcolumn
		for row := Window.Pos.Y; row < Window.Pos.Y+Window.Size.Y; row++ {
			ScreenInstance.Char(tcell.RuneVLine, V2(Window.Pos.X, row), tcell.StyleDefault)
			ScreenInstance.Char(tcell.RuneVLine, V2(Window.Pos.X+Window.Size.X, row), tcell.StyleDefault)
		}
		// draw the top & bottom row
		for col := Window.Pos.X; col < Window.Pos.X+Window.Size.X; col++ {
			ScreenInstance.Char(tcell.RuneHLine, V2(col, Window.Pos.Y), tcell.StyleDefault)
			ScreenInstance.Char(tcell.RuneHLine, V2(col, Window.Pos.Y+Window.Size.Y), tcell.StyleDefault)
		}
		ScreenInstance.Char(tcell.RuneULCorner, V2(Window.Pos.X, Window.Pos.Y), tcell.StyleDefault)
		ScreenInstance.Char('x', V2(Window.Pos.X+Window.Size.X, Window.Pos.Y), tcell.StyleDefault)
		ScreenInstance.Char(tcell.RuneLLCorner, V2(Window.Pos.X, Window.Pos.Y+Window.Size.Y), tcell.StyleDefault)
		ScreenInstance.Char(tcell.RuneLRCorner, V2(Window.Pos.X+Window.Pos.X, Window.Pos.Y+Window.Pos.Y), tcell.StyleDefault)
		ScreenInstance.Text(Window.Title, Window.Pos, tcell.StyleDefault)
	}
}