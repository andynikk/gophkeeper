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

	user.HashPassword = cryptography.HashSHA256(user.Password, dbc.Cfg.Key)
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

	_, err = conn.Exec(ctx, constants.QueryDelOnePairsTemplate, plp.User, plp.Uid)
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

	_, err = conn.Exec(ctx, constants.QueryDelOneBankCardTemplate, td.User, td.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) BankCard(bc *BankCard) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := bc.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		if err = bc.Update(ctx, conn); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}

	if err = bc.Insert(ctx, conn); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) DelBankCard(bc *BankCard) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOneBankCardTemplate, bc.User, bc.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) BinaryData(bd *BinaryData) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := bd.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		if err = bd.Update(ctx, conn); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}
	if err = bd.Insert(ctx, conn); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) DelBinaryData(bd *BinaryData) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	_, err = conn.Exec(ctx, constants.QueryDelOneBinaryDataTemplate, bd.User, bd.Uid)
	if err != nil {
		tx.Rollback(ctx)
		return errs.InvalidFormat
	}
	_, err = conn.Exec(ctx, constants.QueryDelPortionsBinaryData, bd.Uid)
	if err != nil {
		tx.Rollback(ctx)
		return errs.InvalidFormat
	}

	if err := tx.Commit(ctx); err != nil {
		constants.Logger.ErrorLog(err)
		return errs.ErrErrorServer
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

	user.HashPassword = cryptography.HashSHA256(user.Password, dbc.Cfg.Key)
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

		err = rows.Scan(&plp.User, &plp.Uid, &plp.TypePairs, &plp.Name, &plp.Password)
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

func (dbc *DBConnector) SelectBinaryData(ctx context.Context) ([]BinaryData, error) {

	user := ctx.Value(KeyContext("user"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectBinaryData, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrBd []BinaryData
	for rows.Next() {
		var bd BinaryData

		err = rows.Scan(&bd.User, &bd.Uid, &bd.Name, &bd.Expansion, &bd.Size, &bd.Patch)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrBd = append(arrBd, bd)
	}

	return arrBd, nil
}

func (dbc *DBConnector) SelectBankCard(ctx context.Context) ([]BankCard, error) {

	user := ctx.Value(KeyContext("user"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectBankCard, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrBc []BankCard
	for rows.Next() {
		var bc BankCard

		err = rows.Scan(&bc.User, &bc.Uid, &bc.Number, &bc.Cvc)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrBc = append(arrBc, bc)
	}

	return arrBc, nil
}

func (dbc *DBConnector) SelectPortionBinaryData(ctx context.Context) ([]PortionBinaryData, error) {

	uid := ctx.Value(KeyContext("uid"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectPortionsBinaryData, uid)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrPbd []PortionBinaryData
	for rows.Next() {
		var pbd PortionBinaryData

		err = rows.Scan(&pbd.Uid, &pbd.Portion, &pbd.Body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrPbd = append(arrPbd, pbd)
	}

	return arrPbd, nil
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
									OWNER to postgres;`)
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
									"Password" character varying(150) COLLATE pg_catalog."default",
									"UID" character varying(36) COLLATE pg_catalog."default" NOT NULL
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

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."Files"
								(
									"User" character varying(150) COLLATE pg_catalog."default",
									"UID" character varying(36) COLLATE pg_catalog."default",
									"Portion" integer,
									"Name" character varying(150) COLLATE pg_catalog."default",
									"Expansion" character varying(50) COLLATE pg_catalog."default",
									"Body" text COLLATE pg_catalog."default",
									"Patch" character varying(1000) COLLATE pg_catalog."default",
									"Size" character varying COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."Files"
    								OWNER to postgres;`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."PortionsFiles"
								(
									"UID" character varying(36) COLLATE pg_catalog."default",
									"Portion" integer,
									"Body" text COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."PortionsFiles"
									OWNER to postgres;`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."BankCards"
								(
									"User" character varying(150) COLLATE pg_catalog."default",
									"UID" character varying(36) COLLATE pg_catalog."default",
									"Number" character varying COLLATE pg_catalog."default",
									"Date" timestamp with time zone,
									"Cvc" character varying COLLATE pg_catalog."default"
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."BankCards"
									OWNER to postgres;`)
	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return
	}
}
