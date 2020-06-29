package src

import (
	"fmt"
	"github.com/gdamore/tcell"
)

const (
	SCOUT_UNIT = 0x0
)

type Unit interface {
	Update()
	Draw()
	InitUI()
}

type BaseUnit struct {
	Name        string
	Health      float64
	// how many tiles we move per update
	MoveSpeed   float64
	Pos, Target Vec

	DrawManager *DrawManager
}

func NewBaseUnit() *BaseUnit{
	u := &BaseUnit{
		Name:      "base",
		Health:    0,
		MoveSpeed: 1,
		Pos:       V2(15,15),
		Target:    V2(25,25),
		DrawManager: NewDrawManager(),
	}
	u.InitUI()
	return u
}

func (BaseUnit *BaseUnit) InitUI(){
	// add the window for clicking on the settlement
	w := NewWindow(BaseUnit.Name+":window", false, BaseUnit.Name, V2i(10, 10), V2i(20, 20), SCREEN_VIEW)
	w.AddText(NewText("position", true, fmt.Sprintf("pos: %f",BaseUnit.Pos), V2(0, 0), nil, tcell.StyleDefault, SCREEN_VIEW))
	w.AddText(NewText("target", true, fmt.Sprintf("tar: %f",BaseUnit.Target), V2(0, 0), nil, tcell.StyleDefault, SCREEN_VIEW))
	BaseUnit.DrawManager.AddUI(w.ID, w)

	// add the label for the settlement name
	l := NewText(BaseUnit.Name+":label", true, "u", BaseUnit.Pos, func() {
		w.Enable(true)
	}, tcell.StyleDefault.Background(tcell.ColorBlue), WORLD_VIEW)
	BaseUnit.DrawManager.AddUI(l.ID, l)
}

func (BaseUnit *BaseUnit) Update(){
	// get the direction between the 2 vectors
	dir := BaseUnit.Target.Sub(BaseUnit.Pos)
	dir = dir.Normalize().Round()
	BaseUnit.Pos = BaseUnit.Pos.Add(dir)
}

func (BaseUnit *BaseUnit) Draw(){
	BaseUnit.DrawManager.GetText(BaseUnit.Name+":label").Pos = BaseUnit.Pos
	BaseUnit.DrawManager.SetWinText(BaseUnit.Name+":window", "position", fmt.Sprintf("pos: %f",BaseUnit.Pos))
	BaseUnit.DrawManager.SetWinText(BaseUnit.Name+":window", "target", fmt.Sprintf("tar: %f",BaseUnit.Target))
	BaseUnit.DrawManager.Draw()
	//ScreenInstance.Char('u', BaseUnit.Pos, tcell.StyleDefault.Foreground(tcell.ColorYellow), WORLD_VIEW)
}