package src

import (
	"github.com/gdamore/tcell"
)

// represents a settlement that belongs to an empire
// implements Entity
type Settlement struct {
	Empire 	   *Empire
	Name       string
	Population float64
	Pos		   Vec

	// UI
	NameLabel  *Text
	Window     *Window
}

func NewSettlement(empire *Empire, name string, pos Vec) *Settlement{
	s := &Settlement{
		Empire:     empire,
		Name:       name,
		Population: 1.0,
		Pos:        pos,
		NameLabel:  nil,
		Window:     nil,
	}

	s.Window = NewWindow(false, name, V2(10,10), V2(20,20), SCREEN_VIEW)
	s.NameLabel = NewText(true, s.Name, s.Pos.Sub(V2(len(s.Name)/2, 1)), func() {
		s.Window.Enable(true)
	}, tcell.StyleDefault.Background(tcell.ColorBlue), WORLD_VIEW)

	return s
}

func (Settlement *Settlement) Update(){
}

func (Settlement *Settlement) Draw(){
	// draw the settlement symbol
	ScreenInstance.Char('▴', Settlement.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW)
	// draw the name label
	Settlement.NameLabel.Draw()
	Settlement.Window.Draw()
}