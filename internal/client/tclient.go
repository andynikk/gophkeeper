package client

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

// Run Создает CLI-приложение. Создаем и заполняем основную форму с меню
// Описываем события. События - нажате клавишь.
// При старте приложения, создется websocket по адресу "ws://nameserver/socket".
// Каждые две секунды идет опрос сохраненных данных на сервере.
// Данные переносятся на клиент и хранятся в свойстве DataList структуры Client
// На форме отображается и обновляется количество сохраненных записей в базе данных
func (c *Client) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go c.wsData(ctx, cancelFunc)

	c.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {
			c.Pages.SwitchToPage("Menu")
			return nil
		}

		name, _ := c.Pages.GetFrontPage()
		if name != "Menu" {
			return event
		}
		if event.Key() == tcell.KeyCtrlK {
			c.Form.Clear(true)
			c.openEncryptionKeyForms(encryption.KeyRSA{})
			c.Pages.SwitchToPage("KeyRSA")
			return nil
		}

		if c.Name == "" &&
			(event.Rune() == 51 ||
				event.Rune() == 52 ||
				event.Rune() == 53 ||
				event.Rune() == 54 ||
				event.Rune() == 55) {

			c.Pages.SwitchToPage("Menu")
			return nil
		}

		switch event.Rune() {
		case 3: //Ctrl+C
			return nil
		case 48: //0
			c.Application.Stop()
		case 49: //1
			c.Form.Clear(true)
			c.openLoginForm()
			c.Pages.SwitchToPage("Login")
			return nil
		case 50: //2
			c.Form.Clear(true)
			c.openRegisterForms()
			c.Pages.SwitchToPage("Register")
			return nil
		case 51: //3
			c.Form.Clear(true)
			c.openListForms(c.List)
			c.Pages.SwitchToPage("ListData")
			return nil
		case 52: //4
			c.Form.Clear(true)
			c.openPairsLoginPasswordForms(postgresql.PairsLoginPassword{})
			c.Pages.SwitchToPage("PairsLoginPassword")
			return nil
		case 53: //5
			c.Form.Clear(true)
			c.openTextDataForms(postgresql.TextData{})
			c.Pages.SwitchToPage("TextData")
			return nil
		case 54: //6
			c.Form.Clear(true)
			c.openBinaryDataForms(postgresql.BinaryData{})
			c.Pages.SwitchToPage("BinaryData")
			return nil
		case 55: //7
			c.Form.Clear(true)
			c.openBankCardForms(postgresql.BankCard{})
			c.Pages.SwitchToPage("BankCard")
			return nil
		}
		return event
	})

	c.TextView.SetBorder(true)
	c.TextView.SetBorderColor(constants.DefaultColorClient)

	c.TextView.SetTitle(" G O P H K E E P E R ")
	c.TextView.SetTitleColor(constants.DefaultColorClient)

	c.Pages.AddPage("Menu", c.TextView, true, true)
	c.Pages.AddPage("Login", c.Form, true, false)
	c.Pages.AddPage("Register", c.Form, true, false)
	c.Pages.AddPage("PairsLoginPassword", c.Form, true, false)
	c.Pages.AddPage("TextData", c.Form, true, false)
	c.Pages.AddPage("BinaryData", c.Form, true, false)
	c.Pages.AddPage("BankCard", c.Form, true, false)
	c.Pages.AddPage("ListData", c.List, true, false)
	c.Pages.AddPage("KeyRSA", c.Form, true, false)
	c.Pages.AddPage("Comment", c.Form, true, false)

	ctx = context.Background()
	if err := c.Application.SetRoot(c.Pages, true).EnableMouse(true).Sync().Run(); err != nil {
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
	return "USER: " + name + "\n\n" + c.TextDefault + "\n\nRecords counts (" + fmt.Sprintf("%d", i) + ")"
}
