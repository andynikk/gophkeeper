package handlers

import (
	"encoding/json"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/postgresql"
	"gophkeeper/internal/token"
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Server) HandlerNotFound(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, "Page "+r.URL.Path+" not found", http.StatusNotFound)
}

func (srv *Server) HandleFunc(rw http.ResponseWriter, rq *http.Request) {

	//content := srv.StartPage()
	//
	//acceptEncoding := rq.Header.Get("Accept-Encoding")
	//
	//metricsHTML := []byte(content)
	//byteMterics := bytes.NewBuffer(metricsHTML).Bytes()
	//compData, err := compression.Compress(byteMterics)
	//if err != nil {
	//	constants.Logger.ErrorLog(err)
	//}
	//
	//var bodyBate []byte
	//if strings.Contains(acceptEncoding, "gzip") {
	//	rw.Header().Add("Content-Encoding", "gzip")
	//	bodyBate = compData
	//} else {
	//	bodyBate = metricsHTML
	//}
	//
	//rw.Header().Add("Content-Type", "text/html")
	//if _, err := rw.Write(bodyBate); err != nil {
	//	constants.Logger.ErrorLog(err)
	//	return
	//}

	//rw.WriteHeader(http.StatusOK)
}

// POST
func (srv *Server) apiUserRegisterPOST(w http.ResponseWriter, r *http.Request) {

	user := postgresql.User{}
	if err := json.Unmarshal([]byte(r.Header.Get(constants.HeaderMiddlewareBody)), &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenString := ""
	err := srv.DBConnector.NewAccount(&user)
	if err != nil {
		w.Header().Add(constants.HeaderAuthorization, tokenString)
		http.Error(w, err.Error(), HTTPErrors(err))
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

	user := postgresql.User{}

	if err := json.Unmarshal([]byte(r.Header.Get(constants.HeaderMiddlewareBody)), &user); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tokenString := ""
	err := srv.DBConnector.GetAccount(user)
	if err != nil {
		w.Header().Add(constants.HeaderAuthorization, tokenString)
		http.Error(w, err.Error(), HTTPErrors(err))
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

	if err := json.Unmarshal([]byte(r.Header.Get(constants.HeaderMiddlewareBody)), &plp); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelPairsLoginPassword(&plp); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.PairsLoginPassword(&plp); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) apiTextDataPOST(w http.ResponseWriter, r *http.Request) {
	event := mux.Vars(r)["event"]
	td := postgresql.TextData{}

	if err := json.Unmarshal([]byte(r.Header.Get(constants.HeaderMiddlewareBody)), &td); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if event == constants.EventDel.String() {
		if err := srv.DBConnector.DelTextData(&td); err != nil {
			constants.Logger.ErrorLog(err)
			http.Error(w, err.Error(), HTTPErrors(err))
		}
		return
	}

	if err := srv.DBConnector.TextData(&td); err != nil {
		constants.Logger.ErrorLog(err)
		http.Error(w, err.Error(), HTTPErrors(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (srv *Server) Shutdown() {
	constants.Logger.InfoLog("server stopped")
}
