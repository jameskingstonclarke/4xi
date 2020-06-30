package src

import (
	"sync"
)

type GameInterface struct {
	Server *Server
	Client *Client
}

func (G *GameInterface) ProcessServer(){
	G.Server.Init()
	for Running {
		G.Server.Process()
	}
	G.Server.Close()
	WaitGroup.Done()
}

func (G *GameInterface) ProcessClient(){
	G.Client.Init()
	for Running {
		G.Client.Process()
	}
	G.Client.Close()
	WaitGroup.Done()
}

var (
	Running = true
	WaitGroup sync.WaitGroup
)

func Run(){
	InitLogs()
	g := &GameInterface{}
	g.Server = NewServer()
	g.Client = NewClient()
	WaitGroup.Add(2)
	go g.ProcessServer()
	go g.ProcessClient()
	WaitGroup.Wait()
	CloseLogs()
}