package main

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func (c *Client) Send(g *gocui.Gui, v *gocui.View) error {
	text := v.Buffer()

	if text == "" {
		return nil
	}
	sp := []string{text}
	if text[0] == '&' {
		sp = strings.Split(text, " ")
	}
	switch sp[0] {
	case "&ls":
		if sp[1] == "user\n" {
			data := c.LSUSER()
			c.MarshalEndWrite(data)
		}
		if sp[1] == "room\n" {
			data := c.LSROOM()
			c.MarshalEndWrite(data)
		}
	case "&exit\n":
		data := c.EXIT()
		c.MarshalEndWrite(data)
	case "&ch\n":
		data := c.CNAME(sp[1][:len(sp[1])-1])
		c.MarshalEndWrite(data)
	default:
		data := c.Message(text)
		c.MarshalEndWrite(data)
	}

	g.Update(func(g *gocui.Gui) error {
		w, _ := g.View("messages")
		fmt.Fprint(w, text)
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

// Connect test
func Connect(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("messages")
	g.SetViewOnTop("users")
	g.SetViewOnTop("input")
	g.SetCurrentView("input")
	messagesView, _ := g.View("messages")
	g.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(messagesView, "asdasdasd")
		return nil
	})
	return nil
}

// Layout test
func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if messages, err := g.SetView("messages", 0, 0, maxX-20, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = " messages: "
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 0, maxY-5, maxX-20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView("input")
		input.Title = " send: "
		input.Autoscroll = false
		input.Wrap = true
		input.Editable = true
	}

	if users, err := g.SetView("users", maxX-20, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		users.Title = " users: "
		users.Autoscroll = false
		users.Wrap = true
	}

	// if name, err := g.SetView("name", maxX/2-10, maxY/2-1, maxX/2+10, maxY/2+1); err != nil {
	// 	if err != gocui.ErrUnknownView {
	// 		return err
	// 	}
	// 	g.SetCurrentView("name")
	// 	name.Title = " name: "
	// 	name.Autoscroll = false
	// 	name.Wrap = true
	// 	name.Editable = true
	// }
	return nil
}
