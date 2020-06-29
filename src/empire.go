package src

type Empire struct {
	// TODO these are temporary and used as proof-of-concept
	Money  	    float64
	Happyness   float64

	Settlements []Settlement
	Units 		[]Unit
}