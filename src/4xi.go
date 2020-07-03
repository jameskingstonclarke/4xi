package src

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	Game *GameInterface
	Username = "default"
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

func Singleplayer(){
	WaitGroup.Add(2)
	Game.Server = NewServer()
	Game.Client = NewClient("localhost")
	go Game.ProcessServer()
	go Game.ProcessClient()
}

func Host(){
	WaitGroup.Add(2)
	Game.Server = NewServer()
	Game.Client = NewClient("localhost")
	go Game.ProcessServer()
	go Game.ProcessClient()
}

func Connect(){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter server address:")
	addr, err := reader.ReadString('\n')
	if err != nil{
		CLog(err)
	}
	WaitGroup.Add(1)
	Game.Client = NewClient(strings.TrimSpace(addr))
	go Game.ProcessClient()
}

func Run(){

	InitLogs()
	Game = &GameInterface{}



	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter username:")
	username, err := reader.ReadString('\n')
	if err != nil{
		CLog(err)
	}
	Username = strings.TrimSpace(username)

	fmt.Println("singleplayer [0], host game [1], connect to game [2]:")
	choice, err := reader.ReadString('\n')
	if err != nil{
		CLog(err)
	}
	switch strings.TrimSpace(choice){
	case "0":
		Singleplayer()
	case "1":
		Host()
	case "2":
		Connect()
	}

	WaitGroup.Wait()
	CloseLogs()
}