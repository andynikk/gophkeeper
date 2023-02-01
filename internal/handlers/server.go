package handlers

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/midware"
	"gophkeeper/internal/postgresql"
)

type Server struct {
	*mux.Router
	*postgresql.DBConnector
	*environment.ServerConfig
}

func NewServer() *Server {
	srv := &Server{}

	srv.initConfig()
	srv.initDataBase()
	srv.initScoringSystem()
	srv.initRouters()

	return srv
}

// Run Запуск сервера
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
		srv.wsData(conn)
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
	r.Handle("/api/user/pairs/{event}", midware.IsAuthorized(srv.apiPairsLoginPasswordPOST)).Methods("POST")
	r.Handle("/api/user/text/{event}", midware.IsAuthorized(srv.apiTextDataPOST)).Methods("POST")
	r.Handle("/api/user/binary/{event}", midware.IsAuthorized(srv.apiBinaryPOST)).Methods("POST")
	r.Handle("/api/user/card/{event}", midware.IsAuthorized(srv.apiBankCardPOST)).Methods("POST")

	//POST Handle Func
	r.HandleFunc("/api/user/register", srv.apiUserRegisterPOST).Methods("POST")
	r.HandleFunc("/api/user/login", srv.apiUserLoginPOST).Methods("POST")

	r.HandleFunc("/", srv.handleFunc).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(srv.handlerNotFound)
	srv.Router = r
}

func (srv *Server) initDataBase() {
	dbc, err := postgresql.NewDBConnector(&srv.DBConfig)
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
