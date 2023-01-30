package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/postgresql"
)

func (srv *Server) wsData(conn *websocket.Conn) {
	b := []byte("")
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

		ctx := context.Background()
		ctxWV := context.WithValue(ctx, postgresql.KeyContext("user"), string(messageContent))

		arrPlp, err := srv.DBConnector.SelectPairsLoginPassword(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		plpwt := postgresql.PairsLoginPasswordWithType{Type: constants.TypePairsLoginPassword.String(),
			PairsLoginPassword: arrPlp}
		msg, err := json.MarshalIndent(&plpwt, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(1, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}

		arrTd, err := srv.DBConnector.SelectTextData(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		tdwt := postgresql.TextDataWithType{Type: constants.TypeTextData.String(),
			TextData: arrTd}
		msg, err = json.MarshalIndent(&tdwt, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(1, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}

		arrBd, err := srv.DBConnector.SelectBinaryData(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		bdwt := postgresql.BinaryDataWithType{Type: constants.TypeBinaryData.String(),
			BinaryData: arrBd}
		msg, err = json.MarshalIndent(&bdwt, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(1, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}

		arrBc, err := srv.DBConnector.SelectBankCard(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		bcwt := postgresql.BankCardWithType{Type: constants.TypeBankCardData.String(),
			BankCard: arrBc}
		msg, err = json.MarshalIndent(&bcwt, "", " ")
		msg, err = compression.Compress(msg)
		if err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(1, msg); err != nil {
			constants.Logger.ErrorLog(err)
		}

		if err = conn.WriteMessage(0, b); err != nil {
			constants.Logger.ErrorLog(err)
		}

	}
}

func (srv *Server) wsDownloadBinaryData(conn *websocket.Conn, r *http.Request) {

	ctx := context.Background()
	ctxWV := context.WithValue(ctx, postgresql.KeyContext("uid"), r.Header.Get("UID"))

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

		pbd := postgresql.PortionBinaryData{}
		if err = json.Unmarshal(messageContent, &pbd); err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		ctx := context.Background()
		connDB, err := srv.Pool.Acquire(ctx)
		if err != nil {
			return
		}
		if err = pbd.Insert(ctx, connDB); err != nil {
			return
		}
		connDB.Release()
	}
}