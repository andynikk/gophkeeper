package client

import (
	"gophkeeper/internal/environment"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/postgresql"
)

type AuthorizedUser struct {
	postgresql.User
	Token string
}

type Client struct {
	Config *environment.ClientConfig
	AuthorizedUser
	DataList handlers.MapResponse
}

func NewClient() *Client {
	config := environment.InitConfigAgent()
	c := Client{
		Config:         config,
		AuthorizedUser: AuthorizedUser{},
		DataList:       handlers.MapResponse{},
	}

	return &c
}
