package src

import "github.com/gdamore/tcell"

// TODO this is to test a UI system
type UISys struct {
	*SystemBase
	UIManager *UIManager
}

// event to create a new window
type NewWinEvent struct {
	EventBase
	ID    string
	Title string
	Text  []string
}

// event to destroy a window
type DestroyWinEvent struct {
	ID string
}

func (U *UISys) Update(){
	U.UIManager.Draw()
}

func (U *UISys) Remove(){}
func (U *UISys) Priority()int{return 0}

func (U *UISys) ListenNewWinEvent(event NewWinEvent){
	w:=U.UIManager.NewWin(event.ID, true, V2i(10,10))
	for i, t := range event.Text{
		w.NewText(event.ID+":text:"+string(i),t, tcell.StyleDefault, func(){})
	}
}

func (U *UISys) ListenDestroyWinEvent(){}