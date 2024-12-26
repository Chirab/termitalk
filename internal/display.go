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

func (d *Display) renderAddFriendButton() *tview.Button {
	addFriendButton := tview.NewButton("ADD USER").SetSelectedFunc(func() {
		form := tview.NewForm().
			AddInputField("username", "", 20, nil, nil).
			AddButton("Add", func() {
				// call api
			}).
			AddButton("Cancel", func() {
				d.app.SetRoot(d.layout, true).SetFocus(d.list) // Return to main layout
			})

		form.SetBorder(true).SetTitle("Add a Friend").SetTitleAlign(tview.AlignCenter)

		// Show the form as the new root
		d.app.SetRoot(form, true).SetFocus(form)
	})
	addFriendButton.SetBorder(true)
	return addFriendButton
}

func (d *Display) renderRemoveFriendButton() *tview.Button {
	removeFriendButton := tview.NewButton("REMOVE USER").SetSelectedFunc(func() {
		form := tview.NewForm().
			AddInputField("username", "", 20, nil, nil).
			AddButton("remove", func() {
				// call api
			}).
			AddButton("Cancel", func() {
				d.app.SetRoot(d.layout, true).SetFocus(d.list) // Return to main layout
			})

		form.SetBorder(true).SetTitle("remove friend").SetTitleAlign(tview.AlignCenter)

		// Show the form as the new root
		d.app.SetRoot(form, true).SetFocus(form)
	})
	removeFriendButton.SetBorder(true)
	return removeFriendButton
}

func (d *Display) renderSettings() *tview.Flex {
	addFriendBtn := d.renderAddFriendButton()
	removeFriendBtn := d.renderRemoveFriendButton()
	settingsViews := tview.NewFlex().SetDirection(tview.FlexRow)
	settingsViews.SetBorder(true).SetTitle("settings")
	settingsViews.AddItem(addFriendBtn, 0, 1, false)
	settingsViews.AddItem(removeFriendBtn, 0, 1, false)
	return settingsViews
}

func (d *Display) renderLeftSideWindows() *tview.Flex {
	leftWindows := tview.NewFlex().SetDirection(tview.FlexRow)
	leftWindows.AddItem(d.list, 0, 3, true).
		AddItem(d.renderSettings(), 0, 1, false)
	return leftWindows
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

	d.list.Box.SetBorder(true).SetTitle("> Friends <")
	for _, name := range fakeNames() {
		currentName := name
		d.list.AddItem(currentName, "", 'd', func() {
			d.chatUi.SetNewChat(currentName)
		})
	}

	d.layout.SetDirection(tview.FlexColumn).
		AddItem(d.renderLeftSideWindows(), 30, 0, true).
		AddItem(d.chatUi.flex, 0, 3, false)

	d.app.EnableMouse(true)
	if err := d.app.SetRoot(d.layout, true).Run(); err != nil {
		panic(err)
	}
}
