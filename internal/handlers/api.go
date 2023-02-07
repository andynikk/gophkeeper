package handlers

import (
	"encoding/json"
	"gophkeeper/internal/constants/errs"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/postgresql"
	"gophkeeper/internal/token"
)

func (srv *Server) handlerNotFound(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, "Page "+r.URL.Path+" not found", http.StatusNotFound)
}

func (srv *Server) handleFunc(rw http.ResponseWriter, rq *http.Request) {

	if _, err := rw.Write([]byte("Start page")); err != nil {
		constants.Logger.ErrorLog(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (srv *Server) apiUserRegisterPOST(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	user := postgresql.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString := ""
	err = srv.DBConnector.NewAccount(&user)
	if err != nil {
		w.Header().Add(constants.HeaderAuthorization, tokenString)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	tc := token.NewClaims(user.Name)
	if tokenString, err = tc.GenerateJWT(); err != nil {
		w.Header().Add(constants.HeaderAuthorization, "")
		http.Error(w, "Ошибка получения токена", http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.HeaderAuthorization, tokenString)
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiUserLoginPOST(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	user := postgresql.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString := ""
	err = srv.DBConnector.GetAccount(user)
	if err != nil {
		w.Header().Add(constants.HeaderAuthorization, tokenString)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	tc := token.NewClaims(user.Name)
	if tokenString, err = tc.GenerateJWT(); err != nil {
		w.Header().Add(constants.HeaderAuthorization, "")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.HeaderAuthorization, tokenString)
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiPairsLoginPasswordPOST(w http.ResponseWriter, r *http.Request) {
	event := mux.Vars(r)["event"]
	plp := postgresql.PairsLoginPassword{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	if err := json.Unmarshal(body, &plp); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelPairsLoginPassword(&plp); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), errs.HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.PairsLoginPassword(&plp); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiTextDataPOST(w http.ResponseWriter, r *http.Request) {
	event := mux.Vars(r)["event"]
	td := postgresql.TextData{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	if err := json.Unmarshal(body, &td); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelTextData(&td); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), errs.HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.TextData(&td); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiBinaryPOST(w http.ResponseWriter, r *http.Request) {
	event := mux.Vars(r)["event"]
	bd := postgresql.BinaryData{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	if err := json.Unmarshal(body, &bd); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelBinaryData(&bd); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), errs.HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.BinaryData(&bd); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiBankCardPOST(w http.ResponseWriter, r *http.Request) {
	event := mux.Vars(r)["event"]
	bc := postgresql.BankCard{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := r.Header.Get("Content-Encoding")
	if strings.Contains(contentEncoding, "gzip") {
		body, err = compression.Decompress(body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
			return
		}
	}

	if err := json.Unmarshal(body, &bc); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelBankCard(&bc); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), errs.HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.BankCard(&bc); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), errs.HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Shutdown функция отключенния сервера
func (srv *Server) Shutdown() {
	constants.Logger.InfoLog("server stopped")
}
