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
	dest      string
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
		title:     title,
		app:       a,
		chat:      chatView,
		input:     input,
		flex:      flex,
		messageCh: make(chan string),
	}
}

func (c *ChatUi) Reset() {
	c.input.SetText("")
	c.chat.SetText("")
	c.title.SetText("Termitalk")
}

func (c *ChatUi) SetNewChat(name string, p *tview.Pages) {
	if name != c.dest {
		c.dest = name
		c.chat.Clear()
	}

	c.title.SetText(c.dest)
	c.input.SetDoneFunc(func(key tcell.Key) {
		if c.dest == "" {
			return
		}
		if key == tcell.KeyEnter {
			message := c.input.GetText()
			if message == "" {
				return
			}
			if message == "quit" {
				c.app.Stop()
			}
			c.input.SetText("")
			c.messageCh <- fmt.Sprintf("(%s): %s", c.dest, message)
		}
	})

	c.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			p.SwitchToPage("list")
			return nil
		}
		return event
	})
	go func() {
		for msg := range c.messageCh {
			c.chat.Write([]byte(msg + "\n"))
			c.chat.ScrollToEnd()
			c.app.Draw()
		}
	}()
}
