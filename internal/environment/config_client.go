package environment

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"gophkeeper/internal/constants"
)

// ClientConfig структура хранения свойств конфигурации клиента
type ClientConfig struct {
	Address   string
	Key       string
	CryptoKey string
}

type clientConfigENV struct {
	Address   string `env:"ADDRESS" envDefault:"localhost:8080"`
	Key       string `env:"KEY"`
	CryptoKey string `env:"CRYPTO_KEY"`
}

// InitConfigAgent Инициализация и заполнения свойств структуры конфигурации клиента
func InitConfigAgent() *ClientConfig {
	configAgent := ClientConfig{}
	configAgent.InitConfigAgentENV()
	configAgent.InitConfigAgentFlag()

	return &configAgent
}

// InitConfigAgentENV Инициализация и заполнения свойств структуры конфигурации клиента из параметров системы
func (c *ClientConfig) InitConfigAgentENV() {

	var cfgENV clientConfigENV
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
	if fileInfo != nil && err == nil {
		res, err := os.ReadFile(patchCryptoKey)
		if err != nil {
			constants.Logger.ErrorLog(err)
			return
		}
		c.CryptoKey = string(res)
	}
}

// InitConfigAgentFlag Инициализация и заполнения свойств структуры конфигурации клиента из ключей запуска программы
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
		if fileInfo != nil && err == nil {
			res, err := os.ReadFile(*cryptoKeyFlag)
			if err != nil {
				constants.Logger.ErrorLog(err)
				return
			}
			c.CryptoKey = string(res)
		}
	}
}
