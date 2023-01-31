package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/environment"
)

type DBConnector struct {
	Pool *pgxpool.Pool
	Cfg  *environment.DBConfig
}

func NewDBConnector(dbCfg *environment.DBConfig) (*DBConnector, error) {

	if dbCfg.DatabaseDsn == "" {
		return nil, errors.New("пустой путь к базе")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	pool, err := pgxpool.Connect(ctx, dbCfg.DatabaseDsn)
	if err != nil {
		cancelFunc = nil
		return nil, err
	}

	dbc := DBConnector{
		pool,
		dbCfg,
	}

	cancelFunc()
	return &dbc, nil
}
