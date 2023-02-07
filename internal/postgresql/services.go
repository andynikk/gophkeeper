package postgresql

import (
	"context"
	"gophkeeper/internal/token"

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

func (dbc *DBConnector) DelAccount(user *User) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	err = user.Delete(ctx, conn)
	if err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

func (dbc *DBConnector) UpdatePairLoginPassword(plp *PairLoginPassword) error {
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

func (dbc *DBConnector) SelectPairLoginPassword(ctx context.Context) ([]PairLoginPassword, error) {

	user := ctx.Value(KeyContext("user"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectPairsTemplate, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrPlp []PairLoginPassword
	for rows.Next() {
		var plp PairLoginPassword

		err = rows.Scan(&plp.User, &plp.Uid, &plp.TypePairs, &plp.Name, &plp.Password)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrPlp = append(arrPlp, plp)
	}

	return arrPlp, nil
}

func (dbc *DBConnector) DelPairLoginPassword(plp *PairLoginPassword) error {
	claims, ok := token.ExtractClaims(plp.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOnePairsTemplate, claims["user"], plp.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) UpdateTextData(td *TextData) error {
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

func (dbc *DBConnector) DelTextData(td *TextData) error {
	claims, ok := token.ExtractClaims(td.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOneTextDataTemplate, claims["user"], td.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) UpdateBankCard(bc *BankCard) error {
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

func (dbc *DBConnector) DelBankCard(bc *BankCard) error {
	claims, ok := token.ExtractClaims(bc.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, constants.QueryDelOneBankCardTemplate, claims["user"], bc.Uid)
	if err != nil {
		return errs.InvalidFormat
	}
	return nil
}

func (dbc *DBConnector) UpdateBinaryData(bd *BinaryData) error {
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

func (dbc *DBConnector) DelBinaryData(bd *BinaryData) error {
	claims, ok := token.ExtractClaims(bd.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	_, err = conn.Exec(ctx, constants.QueryDelOneBinaryDataTemplate, claims["user"], bd.Uid)
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

func CreateModeLDB(Pool *pgxpool.Pool) error {
	ctx := context.Background()
	conn, err := Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	if _, err = Pool.Exec(ctx, `CREATE SCHEMA IF NOT EXISTS gophkeeper`); err != nil {
		constants.Logger.ErrorLog(err)
		return err
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
		return err
	}

	_, err = conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS gophkeeper."PairLoginPassword"
								(
									"User" character varying(150) COLLATE pg_catalog."default" NOT NULL,
									"TypePairs" character varying(150) COLLATE pg_catalog."default",
									"Name" character varying(150) COLLATE pg_catalog."default",
									"Password" character varying(150) COLLATE pg_catalog."default",
									"UID" character varying(36) COLLATE pg_catalog."default" NOT NULL
								)
								
								TABLESPACE pg_default;
								
								ALTER TABLE IF EXISTS gophkeeper."PairLoginPassword"
									OWNER to postgres;`)

	if err != nil {
		constants.Logger.ErrorLog(err)
		conn.Release()
		return err
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
		return err
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
		return err
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
		return err
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
		return err
	}

	return nil
}
