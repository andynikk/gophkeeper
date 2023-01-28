package environment

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type ClientConfigENV struct {
	Address   string `env:"ADDRESS" envDefault:"localhost:8080"`
	Key       string `env:"KEY"`
	CryptoKey string `env:"CRYPTO_KEY"`
}

type ClientConfig struct {
	Address   string
	Key       string
	CryptoKey string
}

type ClientConfigFile struct {
	Address   string `json:"address"`
	CryptoKey string `json:"crypto_key"`
}

func InitConfigAgent() *ClientConfig {
	configAgent := ClientConfig{}
	configAgent.InitConfigAgentENV()
	configAgent.InitConfigAgentFlag()

	return &configAgent
}

func (c *ClientConfig) InitConfigAgentENV() {

	var cfgENV ClientConfigENV
	err := env.Parse(&cfgENV)
	if err != nil {
		log.Fatal(err)
	}

	addressServ := ""
	if _, ok := os.LookupEnv("ADDRESS"); ok {
		addressServ = cfgENV.Address
	}

	keyHash := ""
	if _, ok := os.LookupEnv("KEY"); ok {
		keyHash = cfgENV.Key
	}

	patchCryptoKey := ""
	if _, ok := os.LookupEnv("CRYPTO_KEY"); ok {
		patchCryptoKey = cfgENV.CryptoKey
	}

	c.Address = addressServ
	c.Key = keyHash
	fileInfo, err := os.Stat(patchCryptoKey)
	if fileInfo != nil && err != nil {
		c.CryptoKey = patchCryptoKey
	}
}

func (c *ClientConfig) InitConfigAgentFlag() {

	addressPtr := flag.String("a", "", "имя сервера")
	keyFlag := flag.String("k", "", "ключ хеширования")
	cryptoKeyFlag := flag.String("c", "", "файл с криптоключем")

	flag.Parse()

	if c.Address == "" {
		c.Address = *addressPtr
	}
	if c.Key == "" {
		c.Key = *keyFlag
	}
	if c.CryptoKey == "" {
		fileInfo, err := os.Stat(*cryptoKeyFlag)
		if fileInfo != nil && err != nil {
			c.CryptoKey = *cryptoKeyFlag
		}
	}
}
