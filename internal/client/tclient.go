package client

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/postgresql"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var form = tview.NewForm()

var arrayEvent = []string{
	"(F9)  Main menu",
	"(F2)  Login",
	"(F3)  Register",
	"(F4)  List info user",
	"(F5)  Add login/password pairs",
	"(F6)  Add arbitrary text data",
	"(F7)  Add arbitrary binary data",
	"(F8)  Add bank card details",
	"(F10) To quit",
	"",
	"(Ctrl+K) Create crypto-key"}

var textDefault = strings.Join(arrayEvent, "\n")

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText(textDefault)
var pages = tview.NewPages()

var list = tview.NewList()

func (c *Client) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	go c.websocket(ctx, cancelFunc)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if c.Name == "" &&
			(event.Key() == tcell.KeyF4 ||
				event.Key() == tcell.KeyF5 ||
				event.Key() == tcell.KeyF6 ||
				event.Key() == tcell.KeyF7 ||
				event.Key() == tcell.KeyF8) {

			text.SetText(c.setMainText())
			pages.SwitchToPage("Menu")
			return nil
		}

		switch event.Key() {
		case tcell.KeyF10:
			app.Stop()
		case tcell.KeyF9:
			text.SetText(c.setMainText())
			pages.SwitchToPage("Menu")
		case tcell.KeyF2:
			form.Clear(true)
			c.loginForm()
			pages.SwitchToPage("Login")
		case tcell.KeyF3:
			form.Clear(true)
			c.registerForms()
			pages.SwitchToPage("Register")
		case tcell.KeyF4:
			form.Clear(true)
			c.listForms(list)
			pages.SwitchToPage("ListData")
		case tcell.KeyF5:
			form.Clear(true)
			c.pairsLoginPasswordForms(postgresql.PairsLoginPassword{})
			pages.SwitchToPage("PairsLoginPassword")
		case tcell.KeyF6:
			form.Clear(true)
			c.textDataForms(postgresql.TextData{})
			pages.SwitchToPage("TextData")
		//case tcell.KeyCtrlK:
		case tcell.KeyF11:
			form.Clear(true)
			c.keyRSAForms(encryption.KeyRSA{})
			pages.SwitchToPage("KeyRSA")
		default:
			//log.Println("++++++++++++++++")
		}
		return event
	})

	pages.AddPage("Menu", text, true, true)
	pages.AddPage("Login", form, true, false)
	pages.AddPage("Register", form, true, false)
	pages.AddPage("PairsLoginPassword", form, true, false)
	pages.AddPage("TextData", form, true, false)
	pages.AddPage("ListData", list, true, false)
	pages.AddPage("KeyRSA", form, true, false)

	if err := app.SetRoot(pages, true).EnableMouse(true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

var done chan interface{}

func (c *Client) receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		msgType, messageContent, err := connection.ReadMessage()
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		if msgType == 0 {
			break
		}

		switch msgType {
		case constants.TypePairsLoginPassword.Int():
			var arrPlp []postgresql.PairsLoginPassword
			if err = json.Unmarshal(messageContent, &arrPlp); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []handlers.Response
			for _, v := range arrPlp {
				arrR = append(arrR, handlers.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(),
					TypeResponse:  msgType})
			}
			c.DataList[constants.TypePairsLoginPassword.String()] = arrR
		case constants.TypeTextData.Int():
			var arrTd []postgresql.TextData
			if err = json.Unmarshal(messageContent, &arrTd); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []handlers.Response
			for _, v := range arrTd {
				arrR = append(arrR, handlers.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(),
					TypeResponse:  msgType})
			}
			c.DataList[constants.TypeTextData.String()] = arrR
		}
	}
}

func (c *Client) websocket(ctx context.Context, cancelFunc context.CancelFunc) {
	socketUrl := "ws://localhost:8080" + "/socket"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	go c.receiveHandler(conn)

	saveTicker := time.NewTicker(time.Duration(3) * time.Second)
	for {
		select {
		case <-saveTicker.C:
			if c.User.Name == "" {
				continue
			}

			err = conn.WriteMessage(1, []byte(c.User.Name))
			if err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}

func (c *Client) setMainText() string {
	name := c.Name
	if name == "" {
		name = "Not authorized"
	}
	//if _, err := os.Stat("/path/to/whatever"); err == nil {
	//
	//}
	return "USER: " + name + "\n\n" + textDefault + "\n\n record counts (" + fmt.Sprintf("%d", len(c.DataList)) + ")"
}
