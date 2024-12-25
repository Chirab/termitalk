package app

import (
	"fmt"
)

type App struct {
	authSvc    IAuth
	displaySvc IDisplay
}

type IAuth interface {
	Login(string, string) bool
	IsAuth() bool
}

type IDisplay interface {
	RenderMain()
}

func NewApp(a IAuth, d IDisplay) *App {
	return &App{
		authSvc:    a,
		displaySvc: d,
	}
}

func (a *App) Run() {
	if !a.authSvc.IsAuth() {
		fmt.Println("please auth")
	}

	a.displaySvc.RenderMain()
}
