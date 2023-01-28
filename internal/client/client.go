package client

import (
	"fmt"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/environment"
	"gophkeeper/internal/handlers"
	"gophkeeper/internal/postgresql"
)

type AuthorizedUser struct {
	postgresql.User
	Token string
}

type Client struct {
	Config        *environment.ClientConfig
	KeyEncryption *encryption.KeyEncryption
	AuthorizedUser
	DataList handlers.MapResponse
}

func NewClient() *Client {
	config := environment.InitConfigAgent()
	c := Client{
		Config:         config,
		KeyEncryption:  new(encryption.KeyEncryption),
		AuthorizedUser: AuthorizedUser{},
		DataList:       handlers.MapResponse{},
	}

	fmt.Println("+++++", config.CryptoKey)
	if config.CryptoKey != "" {
		pk, err := encryption.InitPublicKey(config.CryptoKey)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return nil
		}
		c.KeyEncryption = pk
	}

	return &c
}
