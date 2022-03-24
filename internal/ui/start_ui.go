package ui

import (
	"bufio"
	"cocus/internal/service/client"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

//Init the User Interface and all routines
func StartUi(c *client.Client) {

	fmt.Printf("Nickname for Cocus-Chat: ")

	var uname string
	ConsoleReader := bufio.NewReader(os.Stdin)
	uname, _ = ConsoleReader.ReadString('\n')

	s := strings.TrimSpace(uname)

	ui, err := NewUI(s, c)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.SetManagerFunc(ui.Layout)
	if err := ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Fatalln(err)
	}
	if err := ui.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, ui.Delete); err != nil {
		log.Fatalln(err)
	}
	if err := ui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, ui.SendMsg); err != nil {
		log.Fatalln(err)
	}

	ui.InitMessages()

	go ui.ReceiveMsg()
	if err = ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}

}
