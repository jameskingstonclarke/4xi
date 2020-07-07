package src


import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"reflect"
)

type UI interface {
	Draw()
}

type UIManager struct {
	Screen *Screen
	// TODO optimise this by using a uint32 id
	UI  map[string]UI
}

func (U *UIManager) Remove(id string){
	delete(U.UI, id)
}

func NewUIManager(Screen *Screen) *UIManager{
	return &UIManager{Screen: Screen, UI: make(map[string]UI)}
}

func (UIManager *UIManager) Draw(){
	for _, ui := range UIManager.UI{
		ui.Draw()
	}
}

func (UIManager *UIManager) AddUI(id string, ui UI){
	UIManager.UI[id] = ui
}

func (UIManager *UIManager) NewWin(id string, enabled bool, pos Vec) *Window{
	w := &Window{
		UITemplate: UITemplate{UIManager: UIManager, ID: id, Enabled: enabled, View: SCREEN_VIEW, Style: tcell.StyleDefault},
		Title:      id,
		Pos:        pos,
		Size:       Vec{},
		Text:  		nil,
	}
	UIManager.AddUI(w.ID, w)
	return w
}


func (UIManager *UIManager) NewText(id, text string, enabled bool, pos Vec, view uint8, style tcell.Style, callback Callback) *Text{
	t := &Text{
		UITemplate: UITemplate{UIManager: UIManager, ID: id, Enabled: enabled, View: view, Style: style},
		T:          text,
		Pos:        pos,
		Callback:   callback,
	}
	UIManager.AddUI(t.ID, t)
	return t
}

// set the T value of a text UI element
func (UIManager *UIManager) SetText(id, text string){
	// first check base UI elements
	ui, ok := UIManager.UI[id]
	if !ok{
		// check each window
		CLogErr(errors.New("cannot find UI element "+id))
	}
	switch t:=ui.(type){
	case *Text:
		t.T = text
	}
}

// set the T value of a text UI element
func (UIManager *UIManager) GetText(id string) *Text{
	// first check base UI elements
	ui, ok := UIManager.UI[id]
	if !ok{
		// check each window
		CLogErr(errors.New("cannot find UI element "+id))
	}
	switch t:=ui.(type){
	case *Text:
		return t
	}
	return nil
}

// set the T value of a text UI element
func (UIManager *UIManager) SetWinText(winID, id, text string){
	// first check base UI elements
	win, ok := UIManager.UI[winID]
	if !ok{
		// check each window
		CLogErr(errors.New("cannot find UI element "+id))
	}
	switch w:=win.(type){
	case *Window:
		t := w.WinGetElem(id)
		t.T = text
	}
}

// the reason the text is an interface, is becuase it can be any value and thus be updated on the fly
func (Window *Window) NewText(id string, text interface{}, style tcell.Style, callback Callback){
	t:= &Text{
		UITemplate: UITemplate{UIManager: Window.UIManager, ID: id, Enabled: Window.Enabled, View: Window.View, Style: style},
		T:          text,
		Ptr: 		false,
		Pos:        Vec{},
		Callback:   callback,
	}
	if reflect.ValueOf(text).Type().Kind() == reflect.Ptr{
		t.Ptr = true
	}
	Window.Text = append(Window.Text, t)
	Window.UpdateTextPos()
}

func (Window *Window) WinGetElem(id string) *Text{
	for _, elem := range Window.Text{
		if elem.ID == id{
			return elem
		}
	}
	return nil
}

// set the T value of a text UI element
func (UIManager *UIManager) WinRemoveText(id string, textIDs... string){
	ui, ok := UIManager.UI[id]
	if !ok{
		CLogErr(errors.New("cannot find UI element"+id))
	}
	switch w:=ui.(type){
	case *Window:
		w.RemoveText(textIDs...)
	}
}

type UITemplate struct {
	UIManager *UIManager
	ID      string
	Enabled bool
	View    uint8
	Style   tcell.Style
}

func (UITemplate *UITemplate) Enable(enabled bool){
	UITemplate.Enabled = enabled
}

type Callback func()

type Text struct {
	UITemplate
	// the text is an interface so it can be updated on the fly
	T   	 interface{}
	// if T is a pointer to a value
	Ptr      bool
	Pos 	 Vec
	Callback Callback
}

func (Text *Text) Draw(){
	if Text.Enabled {
		pos := Text.Pos
		if Text.View == WORLD_VIEW{
			pos = Text.UIManager.Screen.WorldToScreen(pos)
		}
		val := fmt.Sprintf("%v",Text.T)
		if Text.Ptr{
			val = fmt.Sprintf("%v",reflect.ValueOf(Text.T).Elem())
		}
		if InputBuffer.MousePressed == '1' {
			if InputBuffer.MousePos.X >= pos.X && int(InputBuffer.MousePos.X) <= int(pos.X)+len(val) {
				if InputBuffer.MousePos.Y == pos.Y {
					if Text.Callback != nil{
						Text.Callback()
					}
				}
			}
		}
		Text.UIManager.Screen.Text(val, Text.Pos, Text.Style, Text.View,UI_DEPTH)
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
		t.Pos = V2i(int(Window.Pos.X)+1,int(Window.Pos.Y) + 1 + i)
		i++
	}
}

// update the size of the window based on the text contents
func (Window *Window) UpdateSize(){
	// find the longest text
	var width, height int
	if len(Window.Title) > width{
		width = len(Window.Title)
	}
	for _, t := range Window.Text{

		val := fmt.Sprintf("%v",t.T)
		if t.Ptr{
			val = fmt.Sprintf("%v",reflect.ValueOf(t.T).Elem())
		}

		if len(val) > width{
			width = len(val)
		}
	}
	height = len(Window.Text)
	Window.Size = V2i(width+1, height+1)
}

func (Window *Window) Draw(){
	if Window.Enabled {
		Window.UpdateSize()
		// check if we are pressing the left mouse
		if InputBuffer.MousePressed == '1'{
			// check if we have clicked the 'x'
			if InputBuffer.MousePos.Equals(Window.Pos.Add(V2(Window.Size.X,0))){
				Window.UIManager.Remove(Window.ID)
			}
		}

		// TODO clean up this window dragging code as it is awful
		if !Window.Dragging && InputBuffer.MouseHeld == '1' && (InputBuffer.MousePos.X >= Window.Pos.X && InputBuffer.MousePos.X <= Window.Pos.X+Window.Size.X &&
			(InputBuffer.MousePos.Y == Window.Pos.Y || InputBuffer.MousePos.Y == Window.Pos.Y+Window.Size.Y)) {
			Window.DragOffset = InputBuffer.MousePos.Sub(Window.Pos)
			Window.Dragging = true
		}
		if Window.Dragging{
			if InputBuffer.MouseHeld == '1' {
				// check if we haven't reached the target
				if !Window.Pos.Equals(InputBuffer.MousePos) {
					Window.UpdateTextPos()
					Window.Pos = InputBuffer.MousePos.Sub(Window.DragOffset)
				}
			}else{
				Window.Dragging = false
			}
		}

		// draw the main body
		for col := Window.Pos.X; col < Window.Pos.X+Window.Size.X; col++ {
			for row := Window.Pos.Y; row < Window.Pos.Y+Window.Size.Y; row++ {
				Window.UIManager.Screen.Char(tcell.RuneBlock, V2(col, row), Window.Style, Window.View,UI_DEPTH)
			}
		}
		// draw the left & right column
		for row := Window.Pos.Y; row < Window.Pos.Y+Window.Size.Y; row++ {
			Window.UIManager.Screen.Char(tcell.RuneVLine, V2(Window.Pos.X, row), Window.Style, Window.View,UI_DEPTH)
			Window.UIManager.Screen.Char(tcell.RuneVLine, V2(Window.Pos.X+Window.Size.X, row), Window.Style, Window.View,UI_DEPTH)
		}
		// draw the top & bottom row
		for col := Window.Pos.X; col < Window.Pos.X+Window.Size.X; col++ {
			Window.UIManager.Screen.Char(tcell.RuneHLine, V2(col, Window.Pos.Y), Window.Style, Window.View,UI_DEPTH)
			Window.UIManager.Screen.Char(tcell.RuneHLine, V2(col, Window.Pos.Y+Window.Size.Y), Window.Style, Window.View,UI_DEPTH)
		}
		Window.UIManager.Screen.Char(tcell.RuneULCorner, V2(Window.Pos.X, Window.Pos.Y), Window.Style, SCREEN_VIEW,UI_DEPTH)
		Window.UIManager.Screen.Char('x', V2(Window.Pos.X+Window.Size.X, Window.Pos.Y), Window.Style, SCREEN_VIEW,UI_DEPTH)
		Window.UIManager.Screen.Char(tcell.RuneLLCorner, V2(Window.Pos.X, Window.Pos.Y+Window.Size.Y), Window.Style, Window.View,UI_DEPTH)
		Window.UIManager.Screen.Char(tcell.RuneLRCorner, V2(Window.Pos.X+Window.Size.X, Window.Pos.Y+Window.Size.Y), Window.Style, Window.View,UI_DEPTH)
		Window.UIManager.Screen.Text(Window.Title, Window.Pos, Window.Style, Window.View,UI_DEPTH)

		for _, t := range Window.Text{
			t.Draw()
		}
	}
}