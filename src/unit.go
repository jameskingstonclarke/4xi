package src

const (
	SCOUT_UNIT = 0x0
)

type Unit interface {
	Update()
	Draw()
}

type BaseUnit struct {
	Type        uint32
	Health      float64
	// how many tiles we move per update
	MoveSpeed   float64
	Pos, Target Vec
}

func (BaseUnit *BaseUnit) Update(){}
func (BaseUnit *BaseUnit) Draw(){}