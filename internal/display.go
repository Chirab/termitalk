package app

import "github.com/rivo/tview"

type Display struct {
	box    *tview.Box
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
		box:    tnl.Box,
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
	for _, name := range fakeNames() {
		currentName := name
		d.list.AddItem(currentName, "", 'd', func() {
			d.chatUi.SetName(currentName)
			d.chatUi.SetMessage()
		})
	}

	d.list.Box.SetBorder(true).SetTitle("> Friends <")
	d.layout.SetDirection(tview.FlexColumn).
		AddItem(d.list, 30, 0, true). // List on the left
		AddItem(d.chatUi.flex, 0, 3, false)

	if err := d.app.SetRoot(d.layout, true).SetFocus(d.list).Run(); err != nil {
		panic(err)
	}
}
