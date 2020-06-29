package src

const (
	SCOUT_UNIT = 0x0
)

type Unit interface {

}

type BaseUnit struct {
	Type      uint32
	Health    float64
	MoveSpeed float64
	Pos       Vec
}