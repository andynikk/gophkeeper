package environment

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"gophkeeper/internal/constants"
)

type ServerConfigENV struct {
	Address     string `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseDsn string `env:"DATABASE_URI"`
	Key         string `env:"KEY"`
}

type DBConfig struct {
	DatabaseDsn string
	Key         string
}

type ServerConfig struct {
	Address string
	DBConfig
}

func NewConfigServer() (*ServerConfig, error) {

	addressPtr := flag.String("a", constants.AdressServer, "адрес сервера")
	keyDatabaseDsn := flag.String("d", "", "строка соединения с базой")
	keyFlag := flag.String("k", "", "ключ хеша")
	flag.Parse()

	var cfgENV ServerConfigENV
	err := env.Parse(&cfgENV)
	if err != nil {
		log.Fatal(err)
	}

	addressServer := cfgENV.Address
	if _, ok := os.LookupEnv("ADDRESS"); !ok {
		addressServer = *addressPtr
	}

	databaseDsn := cfgENV.DatabaseDsn
	if _, ok := os.LookupEnv("DATABASE_URI"); !ok {
		databaseDsn = *keyDatabaseDsn
	}

	keyHash := cfgENV.Key
	if _, ok := os.LookupEnv("KEY"); !ok {
		keyHash = *keyFlag
	}
	if keyHash == "" {
		keyHash = string(constants.HashKey[:])
	}

	sc := ServerConfig{
		Address: addressServer,
		DBConfig: DBConfig{
			DatabaseDsn: databaseDsn,
			Key:         keyHash,
		},
	}

	return &sc, err
}
