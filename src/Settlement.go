package src

import "github.com/nsf/termbox-go"

// represents a settlement that belongs to an empire
type Settlement struct {
	Empire 	   *Empire
	Name       string
	Population int
	X, Y	   int
}

func (Settlement *Settlement) Draw(){
	// draw the name of the settlement
	ScreenInstance.Text(Settlement.Name, Settlement.X-(len(Settlement.Name)/2), Settlement.Y-1, termbox.AttrBold | termbox.ColorGreen, 0)
	// draw the settlement symbol
	ScreenInstance.Char('▴', Settlement.X, Settlement.Y, termbox.AttrBold | termbox.ColorCyan, 0)
}