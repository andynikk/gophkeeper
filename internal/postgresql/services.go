package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/cryptography"
)

func (dbc *DBConnector) NewAccount(user *User) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := user.CheckExistence(ctx, conn)
	if err != nil {
		return errs.ErrErrorServer
	}

	if recordExists {
		return errs.ErrLoginBusy
	}

	user.HashPassword = cryptography.HeshSHA256(user.Password, dbc.Cfg.Key)
	if _, err = conn.Exec(ctx, constants.QueryInsertUserTemplate, user.Name, user.HashPassword); err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

func (dbc *DBConnector) PairsLoginPassword(plp *PairsLoginPassword) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := plp.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		if err = plp.Update(ctx, conn); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}

	if err = plp.Insert(ctx, conn); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) DelPairsLoginPassword(plp *PairsLoginPassword) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOnePairsTemplate, plp.User, plp.TypePairs, plp.Name)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) TextData(td *TextData) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := td.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		if err = td.Update(ctx, conn); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}

	if err = td.Insert(ctx, conn); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) DelTextData(td *TextData) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOneTextDataTemplate, td.User, td.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) GetAccount(user User) error {

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	user.HashPassword = cryptography.HeshSHA256(user.Password, dbc.Cfg.Key)
	recordExists, err := user.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		return nil

	}
	conn.Release()

	return errs.ErrInvalidLoginPassword
}

func (dbc *DBConnector) SelectPairsLoginPassword(ctx context.Context) ([]PairsLoginPassword, error) {

	user := ctx.Value(KeyContext("user"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectPairsTemplate, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrPlp []PairsLoginPassword
	for rows.Next() {
		var plp PairsLoginPassword

		err = rows.Scan(&plp.User, &plp.TypePairs, &plp.Name, &plp.Password)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrPlp = append(arrPlp, plp)
	}

	return arrPlp, nil
}

func (dbc *DBConnector) SelectTextData(ctx context.Context) ([]TextData, error) {

	user := ctx.Value(KeyContext("user"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectTextData, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrTd []TextData
	for rows.Next() {
		var td TextData

		err = rows.Scan(&td.User, &td.Uid, &td.Text)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		arrTd = append(arrTd, td)
	}

	return arrTd, nil
}

func CreateModeLDB(Pool *pgxpool.Pool) {
	ctx := context.Background()
	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return
	}

	if _, err = Pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS gophkeeper`); err != nil {
		constants.Logger.ErrorLog(err)
		return
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."Users"
								(
									"User" character varying(150) COLLATE pg_catalog."default",
									"Password" character varying(256) COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."Users"
									OWNER to postgres;;
`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."PairsLoginPassword"
								(
									"User" character varying(150) COLLATE pg_catalog."default" NOT NULL,
									"TypePairs" character varying(150) COLLATE pg_catalog."default",
									"Name" character varying(150) COLLATE pg_catalog."default",
									"Password" character varying(150) COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."PairsLoginPassword"
									OWNER to postgres;`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."Text"
								(
									"User" character varying(150) COLLATE pg_catalog."default",
									"Text" text COLLATE pg_catalog."default",
									"UID" character varying(36) COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."Text"
									OWNER to postgres;`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}

}
