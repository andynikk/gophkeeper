package client

import (
	"context"
	"fmt"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

// Forms структура хранит данные обекты для создания форм CLI-приложения.
type Forms struct {
	TextDefault string
	*tview.Application
	*tview.Form
	*tview.TextView
	*tview.Pages
	*tview.List
}

// InitForms инициализирует и заполняет структуру Forms стандартными значениями
func InitForms() *Forms {
	var arrayEvent = []string{
		"(ESC) Main menu",
		"(1)   Login",
		"(2)   Register",
		"(3)   List info user",
		"(4)   Add login/password pairs",
		"(5)   Add arbitrary text data",
		"(6)   Add arbitrary binary data",
		"(7)   Add bank card details",
		"(0)   To quit",
		"",
		"(Ctrl+K)  Create crypto-key"}

	textDefault := strings.Join(arrayEvent, "\n")

	return &Forms{
		Application: tview.NewApplication(),
		Form:        tview.NewForm(),
		TextView: tview.NewTextView().
			SetTextColor(constants.DefaultColorClient).
			SetText(textDefault),
		TextDefault: textDefault,
		Pages:       tview.NewPages(),
		List:        tview.NewList(),
	}
}

// Run Создает CLI-приложение. Создаем и заполняем основную форму с меню
// Описываем события. События - нажате клавишь.
// При старте приложения, создется websocket по адресу "ws://nameserver/socket".
// Каждые две секунды идет опрос сохраненных данных на сервере.
// Данные переносятся на клиент и хранятся в свойстве DataList структуры Client
// На форме отображается и обновляется количество сохраненных записей в базе данных
func (f *Forms) Run(c *Client) {

	socketUrl := fmt.Sprintf("ws://%s/socket", c.Config.Address)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		constants.Logger.ErrorLog(err)
		fmt.Println("Ошибка соединения с сервером. Повторите попытку позже")
		return
	}

	ctx := context.Background()

	go c.wsDataWrite(ctx, conn)
	go c.wsDataRead(ctx, conn)
	go f.refreshForm(ctx, c)

	f.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if event.Key() == tcell.KeyEscape {
			f.Pages.SwitchToPage(constants.NameMainPage)
			return nil
		}

		name, _ := f.Pages.GetFrontPage()
		if name != constants.NameMainPage {
			return event
		}
		if event.Key() == tcell.KeyCtrlK {
			f.Form.Clear(true)
			f.openEncryptionKeyForms(c, encryption.KeyRSA{})
			f.Pages.SwitchToPage("KeyRSA")
			return nil
		}

		if c.Name == "" &&
			(event.Rune() == constants.Key3 ||
				event.Rune() == constants.Key4 ||
				event.Rune() == constants.Key5 ||
				event.Rune() == constants.Key6 ||
				event.Rune() == constants.Key7) {

			f.Pages.SwitchToPage(constants.NameMainPage)
			return nil
		}

		switch event.Rune() {
		case constants.KeyCtrlC: //Ctrl+C
			return nil
		case constants.Key0: //0
			f.Application.Stop()
		case constants.Key1: //1
			f.Form.Clear(true)
			f.openLoginForm(c)
			f.Pages.SwitchToPage("Login")
			return nil
		case constants.Key2: //2
			f.Form.Clear(true)
			f.openRegisterForms(c)
			f.Pages.SwitchToPage("Register")
			return nil
		case constants.Key3: //3
			f.Form.Clear(true)
			f.openListForms(c)
			f.Pages.SwitchToPage("ListData")
			return nil
		case constants.Key4: //4
			f.Form.Clear(true)
			f.openPairLoginPasswordForms(c, postgresql.PairLoginPassword{})
			f.Pages.SwitchToPage("PairLoginPassword")
			return nil
		case constants.Key5: //5
			f.Form.Clear(true)
			f.openTextDataForms(c, postgresql.TextData{})
			f.Pages.SwitchToPage("TextData")
			return nil
		case constants.Key6: //6
			f.Form.Clear(true)
			f.openBinaryDataForms(c, postgresql.BinaryData{})
			f.Pages.SwitchToPage("BinaryData")
			return nil
		case constants.Key7: //7
			f.Form.Clear(true)
			f.openBankCardForms(c, postgresql.BankCard{})
			f.Pages.SwitchToPage("BankCard")
			return nil
		}
		return event
	})

	f.TextView.SetBorder(true)
	f.TextView.SetBorderColor(constants.DefaultColorClient)

	f.TextView.SetTitle(" G O P H K E E P E R ")
	f.TextView.SetTitleColor(constants.DefaultColorClient)

	f.Pages.AddPage(constants.NameMainPage, f.TextView, true, true)
	f.Pages.AddPage("Login", f.Form, true, false)
	f.Pages.AddPage("Register", f.Form, true, false)
	f.Pages.AddPage("PairLoginPassword", f.Form, true, false)
	f.Pages.AddPage("TextData", f.Form, true, false)
	f.Pages.AddPage("BinaryData", f.Form, true, false)
	f.Pages.AddPage("BankCard", f.Form, true, false)
	f.Pages.AddPage("ListData", f.List, true, false)
	f.Pages.AddPage("KeyRSA", f.Form, true, false)
	f.Pages.AddPage("Comment", f.Form, true, false)

	if err := f.Application.SetRoot(f.Pages, true).EnableMouse(true).Sync().Run(); err != nil {
		panic(err)
	}
}

// setMainText устанавливает текст основного окна программы
func (f *Forms) setMainText(c *Client) string {
	name := c.Name
	if name == "" {
		name = "Not authorized"
	}
	i := 0
	for _, v := range c.DataList {
		i += len(v)
	}
	return fmt.Sprintf("USER: %s\n\n%s\n\nRecords counts (%d)", name, f.TextDefault, i)
}

// refreshForm горутина которая обновляет текст основного окна программы.
// отображает пользователя и количество записей, хранящихся в БД
func (f *Forms) refreshForm(ctx context.Context, c *Client) {
	ticker := time.NewTicker(time.Second / 2)
	for {
		select {
		case <-ticker.C:
			namePages, _ := f.Pages.GetFrontPage()
			if namePages == constants.NameMainPage {
				newMainText := f.setMainText(c)
				mainText := f.TextView.GetText(true)
				if newMainText != mainText {
					f.TextView.SetText(f.setMainText(c))
					f.Application.ForceDraw()
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
