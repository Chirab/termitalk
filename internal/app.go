package app

import (
	"github.com/rivo/tview"
)

type App struct {
	authSvc    IAuth
	displaySvc IDisplay

	app *tview.Application
}

type IAuth interface {
	RenderAuth()
	IsAuth() bool
}

type IDisplay interface {
	RenderMain()
}

func NewApp() *App {
	app := tview.NewApplication()
	athm := &AuthMemory{}
	return &App{
		authSvc:    NewAuth(app, athm),
		displaySvc: NewDisplay(app),
		app:        app,
	}
}

func (a *App) Run() {
	if !a.authSvc.IsAuth() {
		a.authSvc.RenderAuth()
	} else {

		a.displaySvc.RenderMain()
	}
}
