package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ChatUi struct {
	app       *tview.Application
	chat      *tview.TextView
	input     *tview.InputField
	flex      *tview.Flex
	name      string
	messageCh chan string
	title     *tview.TextView
}

func NewChatUi(a *tview.Application) *ChatUi {
	chatView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)

	input := tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorGreen)

	title := tview.NewTextView().
		SetText("Termitalk").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(title, 1, 0, false).
		AddItem(chatView, 0, 1, false).
		AddItem(input, 3, 0, true)

	return &ChatUi{
		name:      "",
		title:     title,
		app:       a,
		chat:      chatView,
		input:     input,
		flex:      flex,
		messageCh: make(chan string),
	}
}

func (c *ChatUi) SetName(name string) {
	c.name = name
	c.title.SetText(c.name)
}

func (c *ChatUi) SetMessage() {
	c.app.SetFocus(c.flex)
	c.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			message := c.input.GetText()
			if message == "" {
				return
			}
			if message == "quit" {
				c.app.Stop()
			}
			c.input.SetText("")
			c.messageCh <- fmt.Sprintf("(%s): %s", c.name, message)
		}
	})

	go func() {
		for msg := range c.messageCh {
			c.chat.Write([]byte(msg + "\n"))
			c.chat.ScrollToEnd()
			c.app.Draw()
		}
	}()
}
