package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/cryptography"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/postgresql/model"
	"gophkeeper/internal/tests"
	"gophkeeper/internal/token"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/caarlos0/env/v6"
)

type serverConfigENV struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn string `env:"DATABASE_URI"`
	Key         string `env:"KEY"`
}

var srv = &Server{}

func ExampleServer_handlerNotFound() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL+"/not_handler", strings.NewReader(""))
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := fmt.Sprintf("HTTP-Status: %d", resp.StatusCode)
	fmt.Println(msg)

	// Output:
	// HTTP-Status: 404
}

func ExampleServer_apiUserRegisterPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	user := tests.CreateUser("")
	user.HashPassword = cryptography.HashSHA256(user.Password, srv.Key)
	userName := user.Name

	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/user/register", strings.NewReader(string(arrJSON)))
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("User name: %s. HTTP-Status: %d", userName, resp.StatusCode)
	}
	fmt.Println(msg)

	ctx := context.Background()
	conn, err := srv.Pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	pc := model.PgxpoolConn{
		conn,
	}
	ctxVW := context.WithValue(ctx, model.KeyContext("data"), &user)
	err = pc.Delete(ctxVW)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// User name: test. HTTP-Status: 200
}

func ExampleServer_apiUserLoginPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.Background()
	conn, err := srv.Pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()

	user := tests.CreateUser("")
	user.HashPassword = cryptography.HashSHA256(user.Password, srv.Key)
	userName := user.Name

	pc := model.PgxpoolConn{
		conn,
	}
	ctxVW := context.WithValue(ctx, model.KeyContext("data"), &user)
	err = pc.Insert(ctxVW)
	if err != nil {
		return
	}

	arrJSON, err := json.MarshalIndent(user, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/user/login", strings.NewReader(string(arrJSON)))
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("User name: %s. HTTP-Status: %d", userName, resp.StatusCode)
	}
	fmt.Println(msg)

	err = pc.Delete(ctxVW)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// User name: test. HTTP-Status: 200
}

func ExampleServer_apiPairLoginPasswordPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	tc := token.NewClaims("test")
	strToken, _ := tc.GenerateJWT()
	ck := "test crypto key"

	plp := tests.CreatePairLoginPassword(strToken, "", ck)
	uid := plp.Uid

	err := srv.DBConnector.Update(&plp)
	if err != nil {
		return
	}

	arrJSON, err := json.MarshalIndent(plp, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/resource/pairs", strings.NewReader(string(arrJSON)))
	req.Header.Set("Authorization", strToken)
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("UID: %s. HTTP-Status: %d", uid, resp.StatusCode)
	}
	fmt.Println(msg)

	err = srv.DBConnector.Delete(&plp)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// UID: bf340769-687e-485e-968b-976cf12f7b64. HTTP-Status: 200
}

func ExampleServer_apiTextDataPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	tc := token.NewClaims("test")
	strToken, _ := tc.GenerateJWT()
	ck := "test crypto key"

	td := tests.CreateTextData(strToken, "", ck)
	uid := td.Uid

	err := srv.DBConnector.Update(&td)
	if err != nil {
		return
	}

	arrJSON, err := json.MarshalIndent(td, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/resource/text", strings.NewReader(string(arrJSON)))
	req.Header.Set("Authorization", strToken)
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("UID: %s. HTTP-Status: %d", uid, resp.StatusCode)
	}
	fmt.Println(msg)

	err = srv.DBConnector.Delete(&td)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// UID: bf340769-687e-485e-968b-976cf12f7b64. HTTP-Status: 200
}

func ExampleServer_apiBinaryPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	tc := token.NewClaims("test")
	strToken, _ := tc.GenerateJWT()

	bd := tests.CreateBinaryData(strToken, "")
	uid := bd.Uid

	err := srv.DBConnector.Update(&bd)
	if err != nil {
		return
	}

	arrJSON, err := json.MarshalIndent(bd, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/resource/binary", strings.NewReader(string(arrJSON)))
	req.Header.Set("Authorization", strToken)
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("UID: %s. HTTP-Status: %d", uid, resp.StatusCode)
	}
	fmt.Println(msg)

	err = srv.DBConnector.Delete(&bd)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// UID: bf340769-687e-485e-968b-976cf12f7b64. HTTP-Status: 200
}

func ExampleServer_apiBankCardPOST() {
	r := srv.Router
	ts := httptest.NewServer(r)
	defer ts.Close()

	tc := token.NewClaims("test")
	strToken, _ := tc.GenerateJWT()
	ck := "test crypto key"

	bc := tests.CreateBankCard(strToken, "", ck)
	uid := bc.Uid

	err := srv.DBConnector.Update(&bc)
	if err != nil {
		return
	}

	arrJSON, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", ts.URL+"/api/resource/card", strings.NewReader(string(arrJSON)))
	req.Header.Set("Authorization", strToken)
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msg := ""
	if resp.StatusCode == 200 {
		msg = fmt.Sprintf("UID: %s. HTTP-Status: %d", uid, resp.StatusCode)
	}
	fmt.Println(msg)

	err = srv.DBConnector.Delete(&bc)
	if err != nil {
		constants.Logger.ErrorLog(err)
	}

	// Output:
	// UID: bf340769-687e-485e-968b-976cf12f7b64. HTTP-Status: 200
}

func NewConfigServer() (*environment.ServerConfig, error) {

	var cfgENV serverConfigENV
	err := env.Parse(&cfgENV)
	if err != nil {
		return nil, err
	}

	addressServer := cfgENV.Address
	databaseDsn := cfgENV.DatabaseDsn
	keyHash := cfgENV.Key
	if keyHash == "" {
		keyHash = string(constants.HashKey[:])
	}

	sc := environment.ServerConfig{
		Address: addressServer,
		DBConfig: environment.DBConfig{
			DatabaseDsn: databaseDsn,
			Key:         keyHash,
		},
	}

	return &sc, err
}

func init() {

	srv = &Server{}

	srvConfig, err := NewConfigServer()
	if err != nil {
		constants.Logger.ErrorLog(err)
		return
	}
	srv.ServerConfig = srvConfig

	srv.InitDataBase()
	srv.InitRouters()

	srv.InListUserData = InListUserData{}
}
