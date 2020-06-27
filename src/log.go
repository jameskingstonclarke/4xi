package src

import (
	"fmt"
	"os"
)

func Log(data ...interface{}){
	fmt.Println(data...)
}

func LogErr(err error){
	fmt.Println(err)
	os.Exit(2)
}