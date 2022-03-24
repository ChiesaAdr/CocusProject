package ui

import (
	"cocus/internal/service/client"
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

type UI struct {
	*gocui.Gui
	username   string
	transation uint64
	conn       *client.Client
}

func NewUI(username string, conn *client.Client) (*UI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	ui := &UI{Gui: g, username: username, transation: 0, conn: conn}

	return ui, nil
}

func (ui *UI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if messages, err := g.SetView("messages", 0, 0, maxX-20, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = "messages"
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 0, maxY-5, maxX-20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Title = "send"
		input.Autoscroll = false
		input.Wrap = true
		input.Editable = true
	}

	if users, err := g.SetView("users", maxX-20, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		users.Title = "users online"
		users.Autoscroll = false
		users.Wrap = true
	}

	if help, err := g.SetView("help", maxX-40, 0, maxX-20, maxY-45); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		help.Title = "help"
		help.Editable = true
		fmt.Fprintln(help, "KEYBINDINGS")
		fmt.Fprintln(help, "Enter: Send message")
		fmt.Fprintln(help, "^C: Exit")
		fmt.Fprintln(help, "^D: Delete message")
		fmt.Fprintln(help, "")
		fmt.Fprintln(help, "**Attention**")
		fmt.Fprintln(help, " Update only occurs")
		fmt.Fprintln(help, "  with some event!")
	}
	g.SetCurrentView("input")

	return nil
}

func (ui *UI) InitMessages() {
	ui.transation++
	ui.conn.ClientSend(ui.username + "[" + strconv.FormatUint(ui.transation, 10) + "]" + ": online!\n")
}

func (ui *UI) Delete(g *gocui.Gui, v *gocui.View) error {
	fmt.Fprintln(v, "Enter your message[ID]")
	v.SetCursor(0, 1)
	v.FgColor = gocui.ColorRed
	return nil
}

func (ui *UI) Quit(g *gocui.Gui, v *gocui.View) error {
	//Send message to the server to finish
	//Close port and connection
	ui.conn.SendClose2Server()
	return gocui.ErrQuit
}

func (ui *UI) SendMsg(g *gocui.Gui, v *gocui.View) error {
	_, y := v.Cursor()
	if y == 1 {
		b := v.BufferLines()
		if len(b[1]) == 0 {
			return nil
		}
		v.SetCursor(0, 0)
		v.FgColor = gocui.ColorWhite
		v.Clear()
		//Send message to delete by transation ID
		ui.conn.DeleteMessage2Server(b[1])
		return nil
	}

	if len(v.Buffer()) == 0 {
		v.SetCursor(0, 0)
		v.Clear()
		return nil
	}
	ui.transation++
	ui.conn.ClientSend(ui.username + "[" + strconv.FormatUint(ui.transation, 10) + "]" + ": " + v.Buffer())
	v.SetCursor(0, 0)
	v.Clear()
	return nil
}

func (ui *UI) ReceiveMsg() {
	for {

		user := ui.conn.ClientReceive()

		switch user.UserType {
		case 0:
			view, _ := ui.View("messages")
			view.Clear()
			fmt.Fprint(view, strings.Join(user.UserData, ""))
		case 1:
			viewu, _ := ui.View("users")
			viewu.Clear()
			fmt.Fprint(viewu, strings.Join(user.UserData, ""))
		default:
		}
	}
}
