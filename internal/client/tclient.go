package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

var app = tview.NewApplication()
var form = tview.NewForm()

var arrayEvent = []string{
	"(ESC)     Main menu",
	"(Ctrl+L)  Login",
	"(Ctrl+R)  Register",
	"(Ctrl+D)  Info list user",
	"(Ctrl+P)  Add login/password pairs",
	"(Ctrl+T)  Add arbitrary text data",
	"(Ctrl+F)  Add arbitrary binary data",
	"(Ctrl+B)  Add bank card details",
	"(Ctrl+Q)  To quit",
	"",
	"(Ctrl+K)  Create crypto-key"}
var textDefault = strings.Join(arrayEvent, "\n")

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText(textDefault)
var pages = tview.NewPages()

var list = tview.NewList()

func (c *Client) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go c.wsData(ctx, cancelFunc)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if c.Name == "" &&
			(event.Key() == tcell.KeyCtrlD ||
				event.Key() == tcell.KeyCtrlP ||
				event.Key() == tcell.KeyCtrlT ||
				event.Key() == tcell.KeyCtrlB ||
				event.Key() == tcell.KeyF8) {

			text.SetText(c.setMainText())
			pages.SwitchToPage("Menu")
			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyEscape:
			text.SetText(c.setMainText())
			pages.SwitchToPage("Menu")
		case tcell.KeyCtrlL:
			//case tcell.KeyF2:
			form.Clear(true)
			c.openLoginForm()
			pages.SwitchToPage("Login")
		case tcell.KeyCtrlR:
			form.Clear(true)
			c.openRegisterForms()
			pages.SwitchToPage("Register")
		case tcell.KeyCtrlD:
			form.Clear(true)
			c.openListForms(list)
			pages.SwitchToPage("ListData")
		case tcell.KeyCtrlP:
			form.Clear(true)
			c.openPairsLoginPasswordForms(postgresql.PairsLoginPassword{})
			pages.SwitchToPage("PairsLoginPassword")
		case tcell.KeyCtrlT:
			form.Clear(true)
			c.openTextDataForms(postgresql.TextData{})
			pages.SwitchToPage("TextData")
		case tcell.KeyCtrlF:
			form.Clear(true)
			c.openBinaryDataForms(postgresql.BinaryData{})
			pages.SwitchToPage("BinaryData")
		case tcell.KeyCtrlB:
			form.Clear(true)
			c.openBankCardForms(postgresql.BankCard{})
			pages.SwitchToPage("BankCard")
		case tcell.KeyCtrlK:
			form.Clear(true)
			c.openEncryptionKeyForms(encryption.KeyRSA{})
			pages.SwitchToPage("KeyRSA")
		default:

		}
		return event
	})

	pages.AddPage("Menu", text, true, true)
	pages.AddPage("Login", form, true, false)
	pages.AddPage("Register", form, true, false)
	pages.AddPage("PairsLoginPassword", form, true, false)
	pages.AddPage("TextData", form, true, false)
	pages.AddPage("BinaryData", form, true, false)
	pages.AddPage("BankCard", form, true, false)
	pages.AddPage("ListData", list, true, false)
	pages.AddPage("KeyRSA", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (c *Client) setMainText() string {
	name := c.Name
	if name == "" {
		name = "Not authorized"
	}
	i := 0
	for _, v := range c.DataList {
		i += len(v)
	}
	return "USER: " + name + "\n\n" + textDefault + "\n\n record counts (" + fmt.Sprintf("%d", i) + ")"
}
