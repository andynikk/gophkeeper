package environment

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"gophkeeper/internal/constants"
)

type ServerConfigENV struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

type ServerConfig struct {
	Address string
}

func NewConfigServer() (*ServerConfig, error) {

	addressPtr := flag.String("a", constants.AdressServer, "порт сервера")
	flag.Parse()

	var cfgENV ServerConfigENV
	err := env.Parse(&cfgENV)
	if err != nil {
		log.Fatal(err)
	}

	addresServer := cfgENV.Address
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		addresServer = *addressPtr
	}

	sc := ServerConfig{
		Address: addresServer,
	}
	return &sc, err
}
