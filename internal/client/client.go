package client

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/postgresql"
	"strings"

	"github.com/rivo/tview"
)

// ListUserData Список данных пользователя. Заполняется горутиной, котрая запускается
// go c.wsData(ctx, cancelFunc)
// Обновляется каждые 2 секунды
type ListUserData map[string][]postgresql.DataList

// AuthorizedUser структура хранит данные авторизированного пользователя.
// Свойство User хранит имя в явном виде.
// Свойство Token в виде jwt токена.
type AuthorizedUser struct {
	postgresql.User
	Token string
}

// Forms структура хранит данные обекты для создания форм CLI-приложения.
type Forms struct {
	*tview.Application
	*tview.Form
	*tview.TextView
	TextDefault string
	*tview.Pages
	*tview.List
}

// Client общая структура. Хранит все необходимые данные клиента.
type Client struct {
	Config *environment.ClientConfig
	Forms
	AuthorizedUser
	DataList ListUserData
}

// NewClient Создание и заполнение клиента.
// Функция используется при старте клиента.
func NewClient() *Client {
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

	config := environment.InitConfigAgent()
	c := Client{
		Config:         config,
		AuthorizedUser: AuthorizedUser{},
		DataList:       ListUserData{},
		Forms: Forms{
			Application: tview.NewApplication(),
			Form:        tview.NewForm(),
			TextView: tview.NewTextView().
				SetTextColor(constants.DefaultColorClient).
				SetText(textDefault),
			TextDefault: textDefault,
			Pages:       tview.NewPages(),
			List:        tview.NewList(),
		},
	}

	return &c
}
