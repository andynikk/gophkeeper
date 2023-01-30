package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

var done chan interface{}

type TypeMsg struct {
	Type string
}

func (c *Client) wsDownloadBinaryData(ctx context.Context) {

	if c.User.Name == "" {
		return
	}

	abp := ctx.Value(postgresql.KeyContext("additionalBinaryParameters")).(additionalBinaryParameters)

	socketUrl := fmt.Sprintf("ws://%s/socket_download_file", c.Config.Address)
	h := http.Header{}
	h.Add("UID", abp.uid)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, h)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return
	}
	defer conn.Close()

	newFile, err := os.Create(abp.patch)

	for {
		_, messageContent, err := conn.ReadMessage()
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		messageContent, err = compression.Decompress(messageContent)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}

		pbd := postgresql.PortionBinaryData{}
		if err = json.Unmarshal(messageContent, &pbd); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		encryptionBody := encryption.DecryptString(pbd.Body, c.Config.CryptoKey)
		//encryptionBody := pbd.Body
		if _, err = newFile.WriteAt([]byte(encryptionBody), pbd.Portion); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
	}
}

func (c *Client) wsBinaryData(ctx context.Context) {
	socketUrl := fmt.Sprintf("ws://%s/socket_file", c.Config.Address)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return
	}
	defer conn.Close()

	chanOut := make(chan postgresql.PortionBinaryData)
	abp := ctx.Value(postgresql.KeyContext("additionalBinaryParameters")).(additionalBinaryParameters)
	go c.readFile(abp.patch, chanOut)

	for {
		if c.User.Name == "" {
			continue
		}

		bd, ok := <-chanOut
		if !ok {
			break
		}

		bd.Uid = abp.uid
		msg, err := json.MarshalIndent(bd, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		err = conn.WriteMessage(1, msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
	}
}

func (c *Client) wsData(ctx context.Context, cancelFunc context.CancelFunc) {
	socketUrl := fmt.Sprintf("ws://%s/socket", c.Config.Address)
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	go c.receiveHandlerData(conn)

	saveTicker := time.NewTicker(time.Duration(3) * time.Second)
	for {
		select {
		case <-saveTicker.C:
			if c.User.Name == "" {
				continue
			}

			msg := []byte(c.User.Name)
			msg, err := compression.Compress(msg)
			if err != nil {
				constants.Logger.ErrorLog(err)
			}
			err = conn.WriteMessage(1, msg)
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

func (c *Client) receiveHandlerData(connection *websocket.Conn) {
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

		messageContent, err = compression.Decompress(messageContent)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}

		tm := TypeMsg{}
		if err = json.Unmarshal(messageContent, &tm); err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		switch tm.Type {
		case constants.TypePairsLoginPassword.String():
			var PlpWT postgresql.PairsLoginPasswordWithType
			if err = json.Unmarshal(messageContent, &PlpWT); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []postgresql.Response
			arrPlp := PlpWT.PairsLoginPassword
			for _, v := range arrPlp {
				arrR = append(arrR, postgresql.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
					TypeResponse:  constants.TypePairsLoginPassword.String()})
			}
			c.DataList[constants.TypePairsLoginPassword.String()] = arrR
		case constants.TypeTextData.String():
			var Td postgresql.TextDataWithType
			if err = json.Unmarshal(messageContent, &Td); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []postgresql.Response
			arrTd := Td.TextData
			for _, v := range arrTd {
				arrR = append(arrR, postgresql.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
					TypeResponse:  constants.TypeTextData.String()})
			}
			c.DataList[constants.TypeTextData.String()] = arrR
		case constants.TypeBinaryData.String():
			var Bd postgresql.BinaryDataWithType
			if err = json.Unmarshal(messageContent, &Bd); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []postgresql.Response
			arrBd := Bd.BinaryData
			for _, v := range arrBd {
				arrR = append(arrR, postgresql.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
					TypeResponse:  constants.TypeBinaryData.String()})
			}
			c.DataList[constants.TypeBinaryData.String()] = arrR
		case constants.TypeBankCardData.String():
			var Bc postgresql.BankCardWithType
			if err = json.Unmarshal(messageContent, &Bc); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			var arrR []postgresql.Response
			arrBd := Bc.BankCard
			for _, v := range arrBd {
				arrR = append(arrR, postgresql.Response{MainText: v.GetMainText(),
					SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
					TypeResponse:  constants.TypeBankCardData.String()})
			}
			c.DataList[constants.TypeBankCardData.String()] = arrR
		}
	}
}

func (c *Client) readFile(pathSource string, chanOut chan postgresql.PortionBinaryData) {

	file, err := os.Open(pathSource)
	if err != nil {
		close(chanOut)

		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var pos int64 = 0
	for {

		b := make([]byte, constants.Step)
		_, err := file.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			constants.Logger.ErrorLog(err)
		}

		pbd := postgresql.PortionBinaryData{
			Body:    encryption.EncryptString(string(b), c.Config.CryptoKey),
			Portion: pos,
		}

		chanOut <- pbd
		pos += constants.Step
	}
	close(chanOut)
}
