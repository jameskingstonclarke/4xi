package src

import (
	"fmt"
	"github.com/gdamore/tcell"
)

type Settlement interface{
	Update()
	Draw()
	InitUI()
}

// represents a settlement that belongs to an empire
// implements Entity
type BaseSettlement struct {
	Empire 	   *Empire
	Name       string
	Population float64
	PopGrowth  float64
	Pos		   Vec
	UIManager  *UIManager
}

func NewSettlement(empire *Empire, name string, pos Vec) *BaseSettlement{
	s := &BaseSettlement{
		Empire:     empire,
		Name:       name,
		Population: 1,
		PopGrowth:  0.0001,
		Pos:        pos,
		UIManager: NewUIManager(),
	}
	s.InitUI()
	return s
}

func (Settlement *BaseSettlement) Update(){
	Settlement.Population+=Settlement.PopGrowth
}

func (Settlement *BaseSettlement) InitUI(){
	// add the window for clicking on the settlement
	w := NewWindow(Settlement.Name+":window", false, Settlement.Name, V2(10, 10), V2(20, 20), SCREEN_VIEW)
	w.AddText(NewText(Settlement.Name+":population", true, fmt.Sprintf("population: %f",Settlement.Population), V2(0, 0), nil, tcell.StyleDefault, SCREEN_VIEW))
	Settlement.UIManager.AddUI(w.ID, w)

	// add the label for the settlement name
	l := NewText(Settlement.Name+":label", true, Settlement.Name, Settlement.Pos.Sub(V2(len(Settlement.Name)/2, 1)), func() {
		w.Enable(true)
	}, tcell.StyleDefault.Background(tcell.ColorBlue), WORLD_VIEW)
	Settlement.UIManager.AddUI(l.ID, l)
}

func (Settlement *BaseSettlement) Draw(){
	// TODO add this into the UI manager
	ScreenInstance.Char('▴', Settlement.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW)
	Settlement.UIManager.SetWinText(Settlement.Name+":window", Settlement.Name+":population", fmt.Sprintf("population: %f",Settlement.Population))
	Settlement.UIManager.Draw()
}