package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Display struct {
	app    *tview.Application
	list   *tview.List
	layout *tview.Flex

	chatUi *ChatUi
}

func NewDisplay() *Display {
	tnl := tview.NewList()
	app := tview.NewApplication()
	return &Display{

		app:    app,
		list:   tnl,
		layout: tview.NewFlex(),

		chatUi: NewChatUi(app),
	}
}

func fakeNames() []string {
	return []string{
		"spooneyes",
		"saber9516",
		"théo",
		"antoine",
		"moha",
	}
}

func (d *Display) RenderMain() {

	windows := []tview.Primitive{d.list, d.chatUi.flex}
	currentFocus := 0

	d.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			currentFocus = (currentFocus) % len(windows)
			focusedWindow := windows[currentFocus]
			if focusedWindow == d.chatUi.flex {
				d.chatUi.Reset()
			}
			d.app.SetFocus(focusedWindow)

			return nil
		}
		return event
	})
	for _, name := range fakeNames() {
		currentName := name
		d.list.AddItem(currentName, "", 'd', func() {
			d.chatUi.SetNewChat(currentName)
		})
	}
	d.list.Box.SetBorder(true).SetTitle("> Friends <")
	d.layout.SetDirection(tview.FlexColumn).
		AddItem(d.list, 30, 0, true).
		AddItem(d.chatUi.flex, 0, 3, false)

	if err := d.app.SetRoot(d.layout, true).SetFocus(d.list).Run(); err != nil {
		panic(err)
	}
}
