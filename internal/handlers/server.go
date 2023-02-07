package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/midware"
	"gophkeeper/internal/postgresql"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type InListUserData map[string]postgresql.Appender

type Server struct {
	*mux.Router
	*postgresql.DBConnector
	*environment.ServerConfig

	sync.Mutex
	InListUserData
}

// NewServer создание сервера
func NewServer() *Server {
	srv := &Server{}

	srv.InitConfig()
	srv.InitDataBase()
	srv.InitRouters()

	srv.InListUserData = InListUserData{}

	return srv
}

// Run Запуск сервера
func (srv *Server) Run() {
	ctx := context.Background()
	go srv.SaveDataInDB(ctx)

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

// InitRouters инициализация роутера. Описание middleware.
func (srv *Server) InitRouters() {
	r := mux.NewRouter()

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	r.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}

		srv.wsPingData(conn)
	})

	r.HandleFunc("/socket_file", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		srv.wsBinaryData(conn)
	})

	r.HandleFunc("/socket_download_file", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		srv.wsDownloadBinaryData(conn, r)
	})

	//POST
	r.Handle("/api/resource/pairs", midware.IsAuthorized(srv.apiPairLoginPasswordPOST)).Methods("POST")
	r.Handle("/api/resource/text", midware.IsAuthorized(srv.apiTextDataPOST)).Methods("POST")
	r.Handle("/api/resource/binary", midware.IsAuthorized(srv.apiBinaryPOST)).Methods("POST")
	r.Handle("/api/resource/card", midware.IsAuthorized(srv.apiBankCardPOST)).Methods("POST")

	//POST Handle Func
	r.HandleFunc("/api/user/register", srv.apiUserRegisterPOST).Methods("POST")
	r.HandleFunc("/api/user/login", srv.apiUserLoginPOST).Methods("POST")

	r.HandleFunc("/", srv.handleFunc).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(srv.handlerNotFound)
	srv.Router = r
}

// InitDataBase инициализация свойств базы данных сервера
func (srv *Server) InitDataBase() {
	dbc, err := postgresql.NewDBConnector(&srv.DBConfig)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}
	srv.DBConnector = dbc
	_ = postgresql.CreateModeLDB(srv.Pool)
}

// InitConfig инициализация свойств конфигурации сервера
func (srv *Server) InitConfig() {
	srvConfig, err := environment.NewConfigServer()
	if err != nil {
		log.Fatal(err)
	}
	srv.ServerConfig = srvConfig

}

// SaveDataInDB горутина сохранения данных в БД.
// При переброски данных на сервер данные не записываются сразу в БД.
// Данные сохраняются в хранилище сервера InListUserData.
// И только после этого происходит обход InListUserData, перенос данных в БД.
// Удаление из хранилища
func (srv *Server) SaveDataInDB(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:

			srv.SaveData()

		}
	}
}

// SaveData описание непосредственного сохранения данных в БД.
func (srv *Server) SaveData() {
	srv.Lock()
	defer srv.Unlock()

	for kType, vType := range srv.InListUserData {
		for k, v := range vType {
			switch kType {
			case constants.TypePairLoginPassword.String():
				plp := v.(*postgresql.PairLoginPassword)

				if plp.Event == constants.EventDel.String() {
					if err := srv.DBConnector.DelPairLoginPassword(plp); err != nil {
						constants.Logger.ErrorLog(err)
						continue
					}
					delete(vType, k)
					continue
				}

				if err := srv.DBConnector.UpdatePairLoginPassword(plp); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}
				delete(vType, k)
			case constants.TypeTextData.String():
				td := v.(*postgresql.TextData)

				if td.Event == constants.EventDel.String() {
					if err := srv.DBConnector.DelTextData(td); err != nil {
						constants.Logger.ErrorLog(err)
						continue
					}
					delete(vType, k)
					continue
				}

				if err := srv.DBConnector.UpdateTextData(td); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}
				delete(vType, k)
			case constants.TypeBinaryData.String():
				bd := v.(*postgresql.BinaryData)

				if bd.Event == constants.EventDel.String() {
					if err := srv.DBConnector.DelBinaryData(bd); err != nil {
						constants.Logger.ErrorLog(err)
						continue
					}
					delete(vType, k)
					continue
				}

				if err := srv.DBConnector.UpdateBinaryData(bd); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}
				delete(vType, k)
			case constants.TypeBankCardData.String():
				bc := v.(*postgresql.BankCard)

				if bc.Event == constants.EventDel.String() {
					if err := srv.DBConnector.DelBankCard(bc); err != nil {
						constants.Logger.ErrorLog(err)
						continue
					}
					delete(vType, k)
					continue
				}

				if err := srv.DBConnector.UpdateBankCard(bc); err != nil {
					constants.Logger.ErrorLog(err)
					continue
				}
				delete(vType, k)
			default:

			}
		}
	}
}
