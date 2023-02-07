package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/postgresql"
)

// wsBinaryData содает web socket для переброски файлов с клиента на сервер
// Режет файлы на кусочки равные константе Step.
// Шифрует, упаковывает в gzip и отправляет на сервер с разметкой с какого байта начинается.
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

// wsDataWrite, web socket передает на токен залогинящего, текущего пользователя
// Что бы сервер знал какие данные передавать клиенту.
func (c *Client) wsDataWrite(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second / 2)
	for {
		select {
		case <-ticker.C:
			bMsg := []byte(c.Token)
			err := conn.WriteMessage(websocket.TextMessage, bMsg)
			if err != nil {
				constants.Logger.ErrorLog(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// wsDataRead, web socket передает информацию пользователя с сервера на клиент
func (c *Client) wsDataRead(ctx context.Context, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			_, messageContent, err := conn.ReadMessage()
			if err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			messageContent, err = compression.Decompress(messageContent)
			if err != nil {
				constants.Logger.ErrorLog(err)
			}

			tm := postgresql.TypeMsg{}
			if err = json.Unmarshal(messageContent, &tm); err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}

			switch tm.Type {
			case constants.TypePairLoginPassword.String():
				var PlpWT postgresql.PairLoginPasswordWithType
				if err = json.Unmarshal(messageContent, &PlpWT); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}

				var arrR []postgresql.DataList
				arrPlp := PlpWT.Value
				for _, v := range arrPlp {
					arrR = append(arrR, postgresql.DataList{MainText: v.GetMainText(),
						SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
						TypeResponse:  constants.TypePairLoginPassword.String()})
				}
				c.DataList[constants.TypePairLoginPassword.String()] = arrR
			case constants.TypeTextData.String():
				var Td postgresql.TextDataWithType
				if err = json.Unmarshal(messageContent, &Td); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}

				var arrR []postgresql.DataList
				arrTd := Td.TextData
				for _, v := range arrTd {
					arrR = append(arrR, postgresql.DataList{MainText: v.GetMainText(),
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

				var arrR []postgresql.DataList
				arrBd := Bd.BinaryData
				for _, v := range arrBd {
					arrR = append(arrR, postgresql.DataList{MainText: v.GetMainText(),
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

				var arrR []postgresql.DataList
				arrBd := Bc.BankCard
				for _, v := range arrBd {
					arrR = append(arrR, postgresql.DataList{MainText: v.GetMainText(),
						SecondaryText: v.GetSecondaryText(c.Config.CryptoKey),
						TypeResponse:  constants.TypeBankCardData.String()})
				}
				c.DataList[constants.TypeBankCardData.String()] = arrR
			}
		}
	}
}

// readFile, горутина режет файлы на кусочки равные константе Step.
// через какнал chanOut в web socket функции wsBinaryData
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

// wsDownloadBinaryData, web socket передает файл с сервера на клиент и сохраняет на диске
// Получает файл порциями, распаковывает из gzip, расшифровывает. И складывает в один файл
// орентируясь на метки с какого бачта начинается порция
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
		if _, err = newFile.WriteAt([]byte(encryptionBody), pbd.Portion); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
	}
}
