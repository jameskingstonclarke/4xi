package src

import (
	"errors"
	"github.com/gdamore/tcell"
)

type UI interface {
	Draw()
}

type UIManager struct {
	// TODO optimise this by using a uint32 id
	UI  map[string]UI
}

func NewUIManager() *UIManager{
	return &UIManager{UI: make(map[string]UI)}
}

// set the T value of a text UI element
func (UIManager *UIManager) SetText(id, text string){
	ui, ok := UIManager.UI[id]
	if !ok{
		LogErr(errors.New("cannot find UI element"+id))
	}
	switch t:=ui.(type){
	case *Text:
		t.T = text
	}
}

// set the T value of a text UI element
func (UIManager *UIManager) WinAddText(winID string, text... *Text){
	ui, ok := UIManager.UI[winID]
	if !ok{
		LogErr(errors.New("cannot find UI element"+winID))
	}
	switch w:=ui.(type){
	case *Window:
		w.AddText(text...)
	}
}

// set the T value of a text UI element
func (UIManager *UIManager) WinRemoveText(id string, textIDs... string){
	ui, ok := UIManager.UI[id]
	if !ok{
		LogErr(errors.New("cannot find UI element"+id))
	}
	switch w:=ui.(type){
	case *Window:
		w.RemoveText(textIDs...)
	}
}

type UITemplate struct {
	ID      string
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

func NewText(id string, enabled bool, t string, pos Vec, callback Callback, style tcell.Style, view uint8) *Text{
	return &Text{
		UITemplate: UITemplate{ID: id, Enabled: enabled, View: view},
		T:      t,
		Pos:        pos,
		Callback:       callback,
		Style: style,
	}
}


func (Text *Text) Draw(){
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

		ScreenInstance.Text(Text.T, Text.Pos, tcell.StyleDefault, Text.View)
	}
}

type Window struct {
	UITemplate
	Title string
	Pos, Size  Vec
	Text   []*Text
	Dragging bool
	DragOffset Vec
}

func NewWindow(id string, enabled bool, title string, pos, size Vec, view uint8, text ...*Text) *Window{
	w := &Window{
		UITemplate: UITemplate{ID: id, Enabled: enabled, View: view},
		Title:      title,
		Pos:        pos,
		Size:       size,
		Text:  		nil,
	}
	w.AddText(text...)
	return w
}

func (Window *Window) AddText(text ...*Text){
	for _, t := range text {
		Window.Text = append(Window.Text, t)
	}
	Window.UpdateTextPos()
}

func (Window *Window) RemoveText(textIDs ...string){
	// iterate over each text in the window
	for i, text := range Window.Text {
		// check if we need to remove it
		for _, textID := range textIDs{
			if text.ID == textID{
				Window.Text = append(Window.Text[:i], Window.Text[i+1:]...)
			}
		}
	}
	Window.UpdateTextPos()
}

func (Window *Window) UpdateTextPos(){
	i:=0
	for _, t := range Window.Text {
		// the position of the text is the index in the list, plus the border width of the window and the window pos
		t.Pos.Y = (Window.Pos.Y+1) + i
		t.Pos.X = Window.Pos.X+1
		i++
	}
}

// update the size of the window based on the text contents
func (Window *Window) UpdateSize(){
	// find the longest text
	var width, height int
	for _, t := range Window.Text{
		if len(t.T) > width{
			width = len(t.T)
		}
	}
	height = len(Window.Text)
	Window.Size = V2(width+1, height+1)
}

func (Window *Window) Draw(){
	if Window.Enabled {
		Window.UpdateSize()
		// check if we are pressing the left mouse
		if ScreenInstance.InputBuffer.MousePressed & tcell.Button1 != 0{
			// check if we have clicked the 'x'
			if ScreenInstance.InputBuffer.MousePos.Equals(Window.Pos.Add(V2(Window.Size.X,0))){
				Window.Enable(false)
			}
			// check if we are scrolling the window
			if ScreenInstance.InputBuffer.MousePos.X > Window.Pos.X && ScreenInstance.InputBuffer.MousePos.X < Window.Pos.X+Window.Size.X{
				if ScreenInstance.InputBuffer.MousePos.Y == Window.Pos.Y{
					// update the state of the window to indicate it has started or stopped to be dragged
					Window.Dragging = !Window.Dragging
					if Window.Dragging{
						// get the offset of the mouse position so we can subtract it for smooth dragging
						Window.DragOffset = ScreenInstance.InputBuffer.MousePos.Sub(Window.Pos)
					}
				}
			}
		}
		// if we are dragging, update the window
		if Window.Dragging == true{
			// check if we haven't reached the target
			if !Window.Pos.Equals(ScreenInstance.InputBuffer.MousePos) {

				Window.UpdateTextPos()

				Window.Pos = ScreenInstance.InputBuffer.MousePos.Sub(Window.DragOffset)
			}
		}

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



		for _, t := range Window.Text{
			t.Draw()
		}
	}
}