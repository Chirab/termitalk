package main

import (
	"github.com/Chirab/termitalk/internal"
)

func main() {

	auth := app.NewAuth()
	d := app.NewDisplay()
	a := app.NewApp(auth, d)
	a.Run()
}
