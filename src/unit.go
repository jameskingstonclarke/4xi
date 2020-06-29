package src


type Unit struct {
	*Entity
	*PosComp
	*MovementComp
	*RenderComp
}