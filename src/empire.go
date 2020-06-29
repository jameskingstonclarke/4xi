package src
//package src
//
//type Empire struct {
//	Name   		   string
//	// TODO these are temporary and used as proof-of-concept
//	Money, MPT     float64
//	Happyness, HTP float64
//
//	Settlements    []Settlement
//	Units 		   []Unit
//}
//
//func NewEmpire(name string) *Empire{
//	e := &Empire{
//		Name:        name,
//		Money:       0,
//		MPT:         0,
//		Happyness:   0,
//		HTP:         0,
//		Settlements: nil,
//		Units:       nil,
//	}
//	e.Units = append(e.Units, NewBaseUnit())
//	return e
//}
//
//// calculate updated MPT & HPT
//func (Empire *Empire) Update(){
//	for _, settlement := range Empire.Settlements{
//		settlement.Update()
//	}
//	for _, unit := range Empire.Units{
//		unit.Update()
//	}
//}
//
//func (Empire *Empire) Draw(){
//	for _, settlement := range Empire.Settlements{
//		settlement.Draw()
//	}
//	for _, unit := range Empire.Units{
//		unit.Draw()
//	}
//}