package src

import (
	"bufio"
	"fmt"
	"os"
)

const (
	CLogPath = "./clogs.txt"
	SLogPath = "./slogs.txt"
)

var (
	CLogFile *os.File
	CWriter *bufio.Writer
	SLogFile *os.File
	SWriter *bufio.Writer
)

func InitLogs(){
	cf, err := os.Create(CLogPath)
	if err != nil{
		panic(err)
	}
	sf, err := os.Create(SLogPath)
	if err != nil{
		panic(err)
	}
	CLogFile = cf
	SLogFile = sf
	CWriter = bufio.NewWriter(CLogFile)
	SWriter = bufio.NewWriter(SLogFile)
}

func CLog(data ...interface{}){
	fmt.Fprint(CLogFile, data...)
	fmt.Fprintf(CLogFile, "\n")
}

func SLog(data ...interface{}){
	fmt.Fprint(SLogFile, data...)
	fmt.Fprintf(SLogFile, "\n")
}

func CLogErr(err error){
	fmt.Fprint(CLogFile, err)
	fmt.Fprintf(CLogFile, "\n")
	panic(err)
}

func SLogErr(err error){
	fmt.Fprint(SLogFile, err)
	fmt.Fprintf(SLogFile, "\n")
	panic(err)
}

func CloseLogs(){
	CLogFile.Close()
	CWriter.Flush()
	SLogFile.Close()
	SWriter.Flush()
}