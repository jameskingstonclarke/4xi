package src
//
//import (
//	"fmt"
//	"github.com/gdamore/tcell"
//)
//
//type Settlement interface{
//	Update()
//	Draw()
//	InitUI()
//}
//
//// represents a settlement that belongs to an empire
//// implements Entity
//type BaseSettlement struct {
//	Empire 	   *Empire
//	Name       string
//	Population float64
//	PopGrowth  float64
//	Pos		   Vec
//	DrawManager  *DrawManager
//}
//
//func NewSettlement(empire *Empire, name string, pos Vec) *BaseSettlement{
//	s := &BaseSettlement{
//		Empire:     empire,
//		Name:       name,
//		Population: 1,
//		PopGrowth:  0.0001,
//		Pos:        pos,
//		DrawManager: NewDrawManager(),
//	}
//	s.InitUI()
//	return s
//}
//
//func (S *BaseSettlement) Update(){
//	S.Population+=S.PopGrowth
//}
//
//func (S *BaseSettlement) InitUI(){
//	// add the window for clicking on the settlement
//	w:=S.DrawManager.NewWin(S.Name, false, V2i(10, 10), V2i(20, 20), SCREEN_VIEW)
//	w.NewText("pop", fmt.Sprintf("pop: %f",S.Population), true, V2(0, 0), SCREEN_VIEW, tcell.StyleDefault, nil)
//	S.DrawManager.NewText("label", S.Name, true, S.Pos.Sub(V2i(len(S.Name)/2, 1)), WORLD_VIEW,tcell.StyleDefault.Background(tcell.ColorBlue),  func() {
//		w.Enable(true)
//	})
//}
//
//func (S *BaseSettlement) Draw(){
//	// TODO add this into the UI manager
//	ScreenInstance.Char('▴', S.Pos, tcell.StyleDefault.Foreground(tcell.ColorGreen), WORLD_VIEW, 0)
//	S.DrawManager.SetWinText(S.Name, "pop", fmt.Sprintf("pop: %f",S.Population))
//	S.DrawManager.Draw()
//}