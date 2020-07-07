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
	Text  map[interface{}]func()
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
	// sometimes this doesn't execute due to the nature of the goroutines.
	// we need to have a channel that the command listener goroutine passes messages to,
	// this way the event system is thread safe.
	w:=U.UIManager.NewWin(event.ID, true, InputBuffer.MousePos.Sub(V2i(5,5)))
	i:=0
	for t, callback := range event.Text{
		w.NewText(event.ID+":text:"+string(i), t, tcell.StyleDefault, callback)
		i++
	}
}

func (U *UISys) ListenDestroyWinEvent(){}