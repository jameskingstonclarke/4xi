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
	s.Window = &Window{
		UITemplate: UITemplate{Enabled: false},
		Title: name,
		Pos:  V2(10,10),
		Size: V2(20,20),
	}
	s.NameLabel = &Text{
		UITemplate: UITemplate{Enabled: true},
		T:        s.Name,
		Pos:      s.Pos.Sub(V2(len(s.Name)/2, 1)),
		Callback: func() {
			s.Window.Enable(true)
		},
		Style: tcell.StyleDefault.Background(tcell.ColorBlue),
	}
	return s
}

func (Settlement *Settlement) Update(){
	// Process the name label
	Settlement.NameLabel.Update()
	Settlement.Window.Update()
}

func (Settlement *Settlement) Draw(){
	// draw the settlement symbol
	ScreenInstance.Char('▴', Settlement.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	// draw the name label
	Settlement.NameLabel.Draw()
	Settlement.Window.Draw()
}