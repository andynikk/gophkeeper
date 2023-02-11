package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/postgresql/model"
	"gophkeeper/internal/token"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
)

// wsPingData websocket для отправки данных на клиент по имени
func (srv *Server) wsPingData(conn *websocket.Conn) {

	for {
		_, msgToken, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage() error:", err)
			return
		}

		tkn := string(msgToken)
		if tkn == "" {
			continue
		}

		_, ok := token.ExtractClaims(tkn)
		if !ok {
			continue
		}

		app := model.Appender{}

		ctx := context.Background()
		ctxWV := context.WithValue(ctx, model.KeyContext("user"), tkn)

		arrType := []string{constants.TypePairLoginPassword.String(), constants.TypeTextData.String(),
			constants.TypeBinaryData.String(), constants.TypeBankCardData.String()}

		for _, t := range arrType {
			arr, err := srv.DBConnector.Select(ctxWV, t)
			if err != nil {
				constants.Logger.ErrorLog(err)
				continue
			}
			for _, val := range arr {
				app[fmt.Sprintf("%s:%s", t, uuid.New().String())] = val
			}
		}

		msg, err := json.MarshalIndent(&app, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(2, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}
	}
}

// wsDownloadBinaryData websocket переноса бинарных данных с сервера на клиент.
func (srv *Server) wsDownloadBinaryData(conn *websocket.Conn, r *http.Request) {

	ctx := context.Background()
	ctxWV := context.WithValue(ctx, model.KeyContext("uid"), r.Header.Get("UID"))

	arrPbd, err := srv.DBConnector.SelectPortionBinaryData(ctxWV)
	if err != nil {
		constants.Logger.ErrorLog(err)
		return
	}

	for _, v := range arrPbd {
		msg, err := json.MarshalIndent(&v, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(1, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}
	}

	if err = conn.Close(); err != nil {
		constants.Logger.ErrorLog(err)
	}
}

// wsDownloadBinaryData websocket переноса бинарных данных с клиента на сервер.
func (srv *Server) wsBinaryData(conn *websocket.Conn) {
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

		pbd := model.PortionBinaryData{}
		if err = json.Unmarshal(messageContent, &pbd); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		ctx := context.Background()
		connDB, err := srv.Pool.Acquire(ctx)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		defer connDB.Release()

		ctxVW := context.WithValue(ctx, model.KeyContext("data"), pbd)
		if err = srv.DBConnector.InsertPortionBinaryData(ctxVW); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		connDB.Release()
	}
}
