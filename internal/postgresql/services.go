package postgresql

import (
	"context"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/postgresql/model"

	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
)

type InListUserData interface {
	SetFromInListUserData(model.Appender)
}

type DataList struct {
	TypeResponse  string `json:"type"`
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}

// CreateModeLDB при запуске сервера создает таблицы, если их не находит
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
									"User" character varying(150) COLLATE pg_catalog."default" PRIMARY KEY,
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
									"TypePair" character varying(150) COLLATE pg_catalog."default",
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

type PgxpoolConn struct {
	*pgxpool.Conn
}

// CheckExistence проверяет, существует ли объект в базе
func (conn *PgxpoolConn) CheckExistence(ctx context.Context) (bool, error) {

	valCtx := ctx.Value(model.KeyContext("data"))
	b := valCtx.(model.Updater)
	strQuery, argQuery, err := b.CheckExistence()
	if err != nil {
		return false, err
	}
	rows, err := conn.Query(ctx, strQuery, argQuery.([]interface{})...)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Insert добавляет объекты базы данных
func (conn *PgxpoolConn) Insert(ctx context.Context) error {

	valCtx := ctx.Value(model.KeyContext("data"))
	b := valCtx.(model.Updater)
	strQuery, argQuery, err := b.InstructionsInsert()
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, strQuery, argQuery.([]interface{})...)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Update обновляет объекты базы данных
func (conn *PgxpoolConn) Update(ctx context.Context) error {

	valCtx := ctx.Value(model.KeyContext("data"))
	b := valCtx.(model.Updater)
	strQuery, argQuery, err := b.InstructionsUpdate()
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, strQuery, argQuery.([]interface{})...)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Delete удаляет объекты базы данных
func (conn *PgxpoolConn) Delete(ctx context.Context) error {

	valCtx := ctx.Value(model.KeyContext("data"))
	b := valCtx.(model.Updater)
	arrActionDatabase, err := b.InstructionsDelete()
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return errs.InvalidFormat
	}

	for _, v := range arrActionDatabase {
		_, err = conn.Exec(ctx, v.StrExec, v.Arg...)
		if err != nil {
			_ = tx.Rollback(ctx)
			return errs.InvalidFormat
		}
	}

	if err = tx.Commit(ctx); err != nil {
		constants.Logger.ErrorLog(err)
		return errs.InvalidFormat
	}

	return nil
}

// Select выбирает объекты базы данных
func (conn *PgxpoolConn) Select(ctx context.Context) (model.Appender, error) {

	valCtx := ctx.Value(model.KeyContext("data"))
	v := valCtx.(model.Updater)
	actionDatabase, err := v.InstructionsSelect()
	if err != nil {
		return nil, err
	}

	arg := actionDatabase.Arg
	rows, err := conn.Query(ctx, actionDatabase.StrExec, arg...)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	var appender = model.Appender{}

	for rows.Next() {
		app, err := model.NewAppender(actionDatabase.Type, actionDatabase.User)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		err = rows.Scan(app.ArgOut...)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		app.SetValue(appender)
	}

	return appender, nil
}
