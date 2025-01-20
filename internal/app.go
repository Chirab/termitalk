package app

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
)

type App struct {
	authSvc    IAuth
	displaySvc IDisplay

	authMemory *AuthMemory
	app        *tview.Application
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
	app := tview.NewApplication()
	athm := &AuthMemory{}
	return &App{
		authSvc:    NewAuth(app, athm),
		displaySvc: NewDisplay(app),
		app:        app,
		authMemory: athm,
	}
}

func (a *App) Run() {
	maxRetries := 3
	retries := 0

	for retries < maxRetries {
		auth := make(chan AuthState)
		if !a.authSvc.IsAuth() {
			go a.authSvc.RenderAuth(auth)
		}

		state := <-auth
		fmt.Println("state :", state)
		switch state {
		case Islogged:
			fmt.Println("Logged in")
			if err := a.handleLoggedState(); err != nil {
				fmt.Println("Error in logged state, retrying...", err.Error())
				retries++
				continue
			}
			fmt.Println("my server side token :", a.authMemory.GetServerSideToken())
			//a.displaySvc.RenderMain()
			return
		case AuthFailed:
			fmt.Println("Authentication failed, retrying...")
			retries++
			continue
		default:
			fmt.Println("Unhandled state, exiting...")
			return
		}
	}

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
	a.authMemory.DecodeUserInfo(res)
	return nil
}
