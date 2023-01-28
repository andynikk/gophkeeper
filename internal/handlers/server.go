package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/midware"
	"gophkeeper/internal/postgresql"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Storage map[string]MapResponse
type MapResponse map[string][]Response

type Response struct {
	TypeResponse  int    `json:"type"`
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}

type signalEnd struct {
	End bool `json:"end"`
}

type Server struct {
	*mux.Router
	*postgresql.DBConnector
	*environment.ServerConfig
}

func NewByConfig() *Server {
	srv := &Server{}

	srv.initDataBase()
	srv.initConfig()
	srv.initScoringSystem()
	srv.initRouters()

	return srv
}

func (srv *Server) Run() {
	go func() {
		s := &http.Server{
			Addr:    srv.Address,
			Handler: srv.Router}

		if err := s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-stop
	srv.Shutdown()
}

func (srv *Server) initRouters() {
	r := mux.NewRouter()
	r.Use(midware.GzipMiddlware)

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		srv.websocket(conn)
	})

	//POST
	r.Handle("/api/user/pairs/{event}", midware.IsAuthorized(srv.apiPairsLoginPasswordPOST)).Methods("POST")
	r.Handle("/api/user/text/{event}", midware.IsAuthorized(srv.apiTextDataPOST)).Methods("POST")

	//POST Handle Func
	r.HandleFunc("/api/user/register", srv.apiUserRegisterPOST).Methods("POST")
	r.HandleFunc("/api/user/login", srv.apiUserLoginPOST).Methods("POST")

	r.HandleFunc("/", srv.HandleFunc).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(srv.HandlerNotFound)

	srv.Router = r
}

func (srv *Server) websocket(conn *websocket.Conn) {
	b := []byte("")
	for {
		_, messageContent, err := conn.ReadMessage()
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		ctx := context.Background()

		ctxWV := context.WithValue(ctx, postgresql.KeyContext("user"), string(messageContent))
		arrPlp, err := srv.DBConnector.SelectPairsLoginPassword(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		msg, err := json.MarshalIndent(&arrPlp, "", " ")
		if err = conn.WriteMessage(constants.TypePairsLoginPassword.Int(), msg); err != nil {
			constants.Logger.ErrorLog(err)
		}
		arrTd, err := srv.DBConnector.SelectTextData(ctxWV)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		msg, err = json.MarshalIndent(&arrTd, "", " ")
		if err = conn.WriteMessage(constants.TypeTextData.Int(), msg); err != nil {
			constants.Logger.ErrorLog(err)
		}
		if err = conn.WriteMessage(0, b); err != nil {
			constants.Logger.ErrorLog(err)
		}
	}
}

func (srv *Server) initDataBase() {
	dbc, err := postgresql.NewDBConnector()
	if err != nil {
		constants.Logger.ErrorLog(err)
	}
	srv.DBConnector = dbc
	postgresql.CreateModeLDB(srv.Pool)
}

func (srv *Server) initConfig() {
	srvConfig, err := environment.NewConfigServer()
	if err != nil {
		log.Fatal(err)
	}
	srv.ServerConfig = srvConfig

}

func (srv *Server) initScoringSystem() {
	//if !srv.DemoMode {
	//	return
	//}
	//
	//good := Goods{
	//	"My table",
	//	15,
	//	"%",
	//}
	//srv.AddItemsScoringOrder(&good)
	//
	//good = Goods{
	//	"You table",
	//	10,
	//	"%",
	//}
	//srv.AddItemsScoringOrder(&good)
}

func HTTPErrors(err error) int {

	HTTPAnswer := http.StatusOK

	if errors.Is(err, errs.InvalidFormat) {
		HTTPAnswer = http.StatusBadRequest //400
	} else if errors.Is(err, errs.ErrLoginBusy) {
		HTTPAnswer = http.StatusConflict //409
	} else if errors.Is(err, errs.ErrErrorServer) {
		HTTPAnswer = http.StatusInternalServerError //500
	} else if errors.Is(err, errs.ErrInvalidLoginPassword) {
		HTTPAnswer = http.StatusUnauthorized //401
	} else if errors.Is(err, errs.ErrUserNotAuthenticated) {
		HTTPAnswer = http.StatusUnauthorized //401
	} else if errors.Is(err, errs.ErrAccepted) {
		HTTPAnswer = http.StatusAccepted //202
	} else if errors.Is(err, errs.ErrUploadedAnotherUser) {
		HTTPAnswer = http.StatusConflict //409
	} else if errors.Is(err, errs.ErrInvalidOrderNumber) {
		HTTPAnswer = http.StatusUnprocessableEntity //422
	} else if errors.Is(err, errs.ErrInsufficientFunds) {
		HTTPAnswer = http.StatusPaymentRequired //402
	} else if errors.Is(err, errs.ErrNoContent) {
		HTTPAnswer = http.StatusNoContent //204
	} else if errors.Is(err, errs.ErrConflict) {
		HTTPAnswer = http.StatusConflict //409
	} else if errors.Is(err, errs.ErrTooManyRequests) {
		HTTPAnswer = http.StatusTooManyRequests //429
	} else if errors.Is(err, errs.OrderUpload) {
		HTTPAnswer = http.StatusOK //200
	}
	return HTTPAnswer
}
