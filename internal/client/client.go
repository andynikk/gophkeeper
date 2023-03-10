// Package client: запуск клиенского приложения.
package client

import (
	"gophkeeper/internal/environment"
	"gophkeeper/internal/postgresql"
	"gophkeeper/internal/postgresql/model"
)

// ListUserData Список данных пользователя. Заполняется горутиной, котрая запускается
// go c.wsData(ctx, cancelFunc)
// Обновляется каждые 0.5 секунды
type ListUserData map[string][]postgresql.DataList

type BuildInfo struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

// AuthorizedUser структура хранит данные авторизированного пользователя.
// Свойство User хранит имя в явном виде.
// Свойство Token в виде jwt токена.
type AuthorizedUser struct {
	model.User
	Token string
}

// Client общая структура. Хранит все необходимые данные клиента.
type Client struct {
	Config *environment.ClientConfig
	AuthorizedUser
	DataList ListUserData
	BuildInfo
}

// NewClient Создание и заполнение клиента.
// Функция используется при старте клиента.
func NewClient() *Client {
	config := environment.InitConfigAgent()
	c := Client{
		Config:         config,
		AuthorizedUser: AuthorizedUser{},
		DataList:       ListUserData{},
		BuildInfo:      BuildInfo{},
	}

	return &c
}
