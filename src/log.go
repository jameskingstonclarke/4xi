package src

import (
	"bufio"
	"fmt"
	"os"
)

const (
	LogPath = "C:/Users/44778/OneDrive/Programming/GO/src/4xi/logs/logs.txt"
)

var (
	LogFile *os.File
	Writer *bufio.Writer
)

func InitLogs(){
	f, err := os.Create(LogPath)
	if err != nil{
		panic(err)
	}
	LogFile = f
	Writer = bufio.NewWriter(LogFile)
}

func Log(data ...interface{}){
	fmt.Fprint(LogFile, data...)
}

func LogErr(err error){
	fmt.Fprint(LogFile, err)
	panic(err)
}

func CloseLogs(){
	LogFile.Close()
	Writer.Flush()
}