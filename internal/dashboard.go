package app

import (
	"github.com/rivo/tview"
)

type Display struct {
	app    *tview.Application
	list   *tview.List
	layout *tview.Flex

	chatUi *ChatUi

	pages *tview.Pages
}

func NewDisplay(app *tview.Application) *Display {
	tnl := tview.NewList()
	pages := tview.NewPages()
	return &Display{

		app:    app,
		list:   tnl,
		layout: tview.NewFlex(),

		chatUi: NewChatUi(app),
		pages:  pages,
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

func (d *Display) renderAddFriendForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("username", "", 20, nil, nil).
		AddButton("Add", func() {
			// call api
		})

	form.SetBorder(true).SetTitle("Add a Friend").SetTitleAlign(tview.AlignCenter)
	d.app.SetRoot(form, true).SetFocus(form)

	return form
}

func (d *Display) renderRemoveFriendForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("username", "", 20, nil, nil).
		AddButton("remove", func() {
			// call api
		})

	form.SetBorder(true).SetTitle("remove friend").SetTitleAlign(tview.AlignCenter)
	d.app.SetRoot(form, true).SetFocus(form)
	return form
}

func (d *Display) renderSettings() *tview.Flex {
	addFriendBtn := d.renderAddFriendForm()
	removeFriendBtn := d.renderRemoveFriendForm()
	settingsViews := tview.NewFlex().SetDirection(tview.FlexRow)
	settingsViews.SetBorder(true).SetTitle("SETTINGS")
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

func (d *Display) renderUsername(s string) *tview.TextView {
	username := tview.NewTextView().
		SetText(s).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	return username
}

func (d *Display) renderFriendRequests() *tview.Flex {
	friendRequest := tview.NewFlex().SetDirection(tview.FlexRow)
	friendRequest.AddItem(d.renderUsername("olivier"), 5, 0, true)

	return friendRequest
}

func (d *Display) RenderMain() {
	d.list.Box.SetBorder(true).SetTitle("> PRIVATE MESSAGES <")

	d.pages.AddPage("chat", d.chatUi.flex, true, false)
	for _, name := range fakeNames() {
		currentName := name
		d.list.AddItem(currentName, "", ' ', func() {
			d.chatUi.SetNewChat(currentName, d.pages)
			d.pages.SwitchToPage("chat")

		})
	}

	d.layout.SetDirection(tview.FlexColumn).
		AddItem(d.renderLeftSideWindows(), 0, 1, true).
		AddItem(d.renderFriendRequests(), 0, 1, true)
	d.pages.AddPage("list", d.layout, true, true)
	d.app.EnableMouse(true)
	if err := d.app.SetRoot(d.pages, true).Run(); err != nil {
		panic(err)
	}
}
