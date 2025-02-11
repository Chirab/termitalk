package app

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
)

type App struct {
	authSvc    IAuth
	displaySvc IDisplay
	serverSvc  *Server

	authMemory *Memory
	app        *tview.Application
	cm         chan string
}

type IAuth interface {
	RenderAuth(chan AuthState)
	CheckIsLogged() (string, error)
	IsAuth() bool
	UnsetAppRoot()
}

type IDisplay interface {
	RenderMain()
}

type AuthState int

const (
	Islogged = iota + 1
	AuthFailed
)

func NewApp() *App {
	cm := make(chan string)
	athm := &Memory{}
	athm.InitMetadataStorage()
	srv := NewServer("3434", athm, cm)
	app := tview.NewApplication()

	return &App{
		authSvc:    NewAuth(app, athm),
		displaySvc: NewDisplay(app),
		app:        app,
		authMemory: athm,
		serverSvc:  srv,
		cm:         cm,
	}
}

func (a *App) Run() {
	maxRetries := 3
	retries := 0

	go a.serverSvc.Start()

	for retries < maxRetries {
		auth := make(chan AuthState)
		if !a.authSvc.IsAuth() {
			go a.authSvc.RenderAuth(auth)
		}

		authId := <-a.cm
		a.authMemory.SetAuthId(authId)

		state := <-auth
		switch state {
		case Islogged:
			if err := a.handleLoggedState(); err != nil {
				fmt.Println("Error in logged state, retrying...", err.Error())
				retries++
				continue
			}
			a.serverSvc.Shutdown()
			a.displaySvc.RenderMain()
		case AuthFailed:
			fmt.Println("Authentication failed, retrying...")
			retries++
			continue
		default:
			fmt.Println("Unhandled state, exiting...")
			return
		}
	}

	defer a.serverSvc.Shutdown()

	fmt.Println("Max retries reached, exiting...")

}

func (a *App) handleLoggedState() error {
	res, err := a.authSvc.CheckIsLogged()
	if err != nil {
		return err
	}
	if res == "" {
		return errors.New("Logged in was not successful")
	}
	return nil
}
