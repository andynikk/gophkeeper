package client

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/postgresql"
	"strings"

	"github.com/rivo/tview"
)

type AuthorizedUser struct {
	postgresql.User
	Token string
}

type Forms struct {
	*tview.Application
	*tview.Form
	*tview.TextView
	TextDefault string
	*tview.Pages
	*tview.List
}

type Client struct {
	Config *environment.ClientConfig
	Forms
	AuthorizedUser
	DataList handlers.MapResponse
}

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
		DataList:       handlers.MapResponse{},
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
