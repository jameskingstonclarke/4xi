package src

const (
	CLIENT = 0x0
	HOST   = 0x1
)

var (
	Mode = HOST

)

func Run(){
	g := &Game{}
	g.Init()
	g.Process()
	g.Close()
}