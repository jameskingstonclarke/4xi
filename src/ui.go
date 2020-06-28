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
	View    uint8
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

func NewText(enabled bool, t string, pos Vec, callback Callback, style tcell.Style, view uint8) *Text{
	return &Text{
		UITemplate: UITemplate{Enabled: enabled, View: view},
		T:      t,
		Pos:        pos,
		Callback:       callback,
		Style: style,
	}
}


func (Text *Text) Update(){
	if Text.Enabled {
		pos := Text.Pos
		if Text.View == WORLD_VIEW{
			pos = ScreenInstance.WorldToScreen(pos)
		}
		if ScreenInstance.InputBuffer.MousePressed&tcell.Button1 != 0 {
			if ScreenInstance.InputBuffer.MousePos.X > pos.X && ScreenInstance.InputBuffer.MousePos.X < pos.X+len(Text.T) {
				if ScreenInstance.InputBuffer.MousePos.Y == pos.Y {
					Text.Callback()
				}
			}
		}
	}
}

func (Text *Text) Draw(){
	if Text.Enabled {
		ScreenInstance.Text(Text.T, Text.Pos, tcell.StyleDefault, Text.View)
	}
}

type Window struct {
	UITemplate
	Title string
	Pos, Size  Vec
}

func NewWindow(enabled bool, title string, pos, size Vec, view uint8) *Window{
	return &Window{
		UITemplate: UITemplate{Enabled: enabled, View: view},
		Title:      title,
		Pos:        pos,
		Size:       size,
	}
}

func (Window *Window) Update(){
	if Window.Enabled{
		pos := Window.Pos
		if Window.View == WORLD_VIEW{
			pos = ScreenInstance.WorldToScreen(pos)
		}
		if ScreenInstance.InputBuffer.MousePressed & tcell.Button1 != 0{
			if ScreenInstance.InputBuffer.MousePos.Equals(pos.Add(V2(Window.Size.X,0))){
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
				ScreenInstance.Char(tcell.RuneBlock, V2(col, row), tcell.StyleDefault, Window.View)
			}
		}
		// draw the left & right column
		for row := Window.Pos.Y; row < Window.Pos.Y+Window.Size.Y; row++ {
			ScreenInstance.Char(tcell.RuneVLine, V2(Window.Pos.X, row), tcell.StyleDefault, Window.View)
			ScreenInstance.Char(tcell.RuneVLine, V2(Window.Pos.X+Window.Size.X, row), tcell.StyleDefault, Window.View)
		}
		// draw the top & bottom row
		for col := Window.Pos.X; col < Window.Pos.X+Window.Size.X; col++ {
			ScreenInstance.Char(tcell.RuneHLine, V2(col, Window.Pos.Y), tcell.StyleDefault, Window.View)
			ScreenInstance.Char(tcell.RuneHLine, V2(col, Window.Pos.Y+Window.Size.Y), tcell.StyleDefault, Window.View)
		}
		ScreenInstance.Char(tcell.RuneULCorner, V2(Window.Pos.X, Window.Pos.Y), tcell.StyleDefault, SCREEN_VIEW)
		ScreenInstance.Char('x', V2(Window.Pos.X+Window.Size.X, Window.Pos.Y), tcell.StyleDefault, SCREEN_VIEW)
		ScreenInstance.Char(tcell.RuneLLCorner, V2(Window.Pos.X, Window.Pos.Y+Window.Size.Y), tcell.StyleDefault, Window.View)
		ScreenInstance.Char(tcell.RuneLRCorner, V2(Window.Pos.X+Window.Size.X, Window.Pos.Y+Window.Size.Y), tcell.StyleDefault, Window.View)
		ScreenInstance.Text(Window.Title, Window.Pos, tcell.StyleDefault, Window.View)
	}
}