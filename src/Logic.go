package src

type Logic struct {

}

var (
	LogicInstance *Logic
)

func (Logic *Logic) Process(){
	for Running{

	}
}

func (Logic *Logic) Init(){
}

func (Logic *Logic) Close(){
	WaitGroup.Done()
}