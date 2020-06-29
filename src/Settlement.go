package src

import (
	"github.com/gdamore/tcell"
)

type Settlement interface{}

// represents a settlement that belongs to an empire
// implements Entity
type BaseSettlement struct {
	Empire 	   *Empire
	Name       string
	Population float64
	Pos		   Vec

	// UI
	NameLabel  *Text
	Window     *Window
}

func NewSettlement(empire *Empire, name string, pos Vec) *BaseSettlement{
	s := &BaseSettlement{
		Empire:     empire,
		Name:       name,
		Population: 1.0,
		Pos:        pos,
		NameLabel:  nil,
		Window:     nil,
	}

	s.Window = NewWindow("settlement_window",false, name, V2(10,10), V2(20,20), SCREEN_VIEW)
	s.Window.AddText(NewText("l1",true, "line of text 1", V2(0,0), nil, tcell.StyleDefault, SCREEN_VIEW))
	s.Window.AddText(NewText("l2",true, "line of text 2", V2(0,0), nil, tcell.StyleDefault, SCREEN_VIEW))
	s.Window.AddText(NewText("l3",true, "line of text 3", V2(0,0), nil, tcell.StyleDefault, SCREEN_VIEW))
	s.NameLabel = NewText("settlement_name",true, s.Name, s.Pos.Sub(V2(len(s.Name)/2, 1)), func() {
		s.Window.Enable(true)
	}, tcell.StyleDefault.Background(tcell.ColorBlue), WORLD_VIEW)

	return s
}

func (Settlement *BaseSettlement) Update(){
}

func (Settlement *BaseSettlement) Draw(){
	// draw the settlement symbol
	ScreenInstance.Char('▴', Settlement.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW)
	// draw the name label
	Settlement.NameLabel.Draw()
	Settlement.Window.Draw()
}