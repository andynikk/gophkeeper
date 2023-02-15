// Package handlers: хендлеры сервера. Отработка действи хендлера
package handlers

import (
	"encoding/json"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/postgresql/model"
	"io"
	"net/http"
	"strings"

	"gophkeeper/internal/compression"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/token"
)

// handlerNotFound, хендлер адрес не найден
func (srv *Server) handlerNotFound(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, "Page "+r.URL.Path+" not found", http.StatusNotFound)
}

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFunc(rw http.ResponseWriter, rq *http.Request) {

	if _, err := rw.Write([]byte("Start page")); err != nil {
		constants.Logger.ErrorLog(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// apiUserRegisterPOST хендлер создания пользователя
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

	user := model.User{}
	if err := json.Unmarshal(body, &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString := ""
	user.New = true
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

// apiUserRegisterPOST хендлер входа пользователя в систему
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

	user := model.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenString := ""
	err = srv.DBConnector.CheckAccount(&user)
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

// apiPairLoginPasswordPOST хендлер для работы с данными типа "пары логин/пароль"
func (srv *Server) apiPairLoginPasswordPOST(w http.ResponseWriter, r *http.Request) {

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

	var plp model.PairLoginPassword
	if err = json.Unmarshal(body, &plp); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, "Ошибка распаковки", http.StatusInternalServerError)
		return
	}

	plp.User = r.Header.Get("Authorization")

	srv.Mutex.Lock()
	defer srv.Mutex.Unlock()

	plpInListUserData, ok := srv.InListUserData[constants.TypePairLoginPassword.String()]
	if !ok {
		plpInListUserData = model.Appender{}
	}
	plp.SetFromInListUserData(plpInListUserData)

	srv.InListUserData[constants.TypePairLoginPassword.String()] = plpInListUserData
	w.WriteHeader(http.StatusOK)
}

// apiTextDataPOST хендлер для работы с данными типа "произвольные текстовые данные"
func (srv *Server) apiTextDataPOST(w http.ResponseWriter, r *http.Request) {
	td := model.TextData{}

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

	td.User = r.Header.Get("Authorization")

	srv.Mutex.Lock()
	defer srv.Mutex.Unlock()

	tdInListUserData, ok := srv.InListUserData[constants.TypeTextData.String()]
	if !ok {
		tdInListUserData = model.Appender{}
	}
	td.SetFromInListUserData(tdInListUserData)

	srv.InListUserData[constants.TypeTextData.String()] = tdInListUserData
	w.WriteHeader(http.StatusOK)

}

// apiBinaryPOST хендлер для работы с данными типа "произвольные бинарные данные"
func (srv *Server) apiBinaryPOST(w http.ResponseWriter, r *http.Request) {
	bd := model.BinaryData{}

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

	bd.User = r.Header.Get("Authorization")

	srv.Mutex.Lock()
	defer srv.Mutex.Unlock()

	tdInListUserData, ok := srv.InListUserData[constants.TypeBinaryData.String()]
	if !ok {
		tdInListUserData = model.Appender{}
	}
	bd.SetFromInListUserData(tdInListUserData)

	srv.InListUserData[constants.TypeBinaryData.String()] = tdInListUserData

	w.WriteHeader(http.StatusOK)
}

// apiBankCardPOST хендлер для работы с данными типа "данные банковских карт"
func (srv *Server) apiBankCardPOST(w http.ResponseWriter, r *http.Request) {

	bc := model.BankCard{}

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

	bc.User = r.Header.Get("Authorization")

	srv.Mutex.Lock()
	defer srv.Mutex.Unlock()

	tdInListUserData, ok := srv.InListUserData[constants.TypeBankCardData.String()]
	if !ok {
		tdInListUserData = model.Appender{}
	}
	bc.SetFromInListUserData(tdInListUserData)
	srv.InListUserData[constants.TypeBankCardData.String()] = tdInListUserData

	w.WriteHeader(http.StatusOK)
}

// Shutdown функция, работающая при отключеннии сервера
func (srv *Server) Shutdown() {
	constants.Logger.InfoLog("server stopped")
}
