package postgresql

import (
	"context"
	"gophkeeper/internal/postgresql/model"

	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/cryptography"
)

type AppenderWithType struct {
	Type  string         `json:"type"`
	Event string         `json:"event"`
	Value model.Appender `json:"value"`
}

type KeyMsg struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type TypeMsg struct {
	Type  string
	Token string
}

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

// NewAccount метод для создания нового экаунта из ДБ конектора
// вызывает методы объекта user.
// Проверяет есть ли такой пользователь.
// Если нет, то создает
func (dbc *DBConnector) NewAccount(user *model.User) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	pc := model.PgxpoolConn{conn}
	recordExists, err := pc.CheckExistence(ctxVW)
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

// CheckAccount проверяет, по имени и хешированному паролю существует ли пользователь в базе
func (dbc *DBConnector) CheckAccount(user *model.User) error {

	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	user.HashPassword = cryptography.HashSHA256(user.Password, dbc.Cfg.Key)
	ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	var pc = model.PgxpoolConn{Conn: conn}
	recordExists, err := pc.CheckExistence(ctxVW)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		return nil

	}
	conn.Release()

	return errs.ErrInvalidLoginPassword
}

// DelAccount удаляет пользователя по имени и хешированному паролю
func (dbc *DBConnector) DelAccount(user *model.User) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	ctxVW := context.WithValue(ctx, model.KeyContext("data"), user)
	pc := model.PgxpoolConn{conn}

	err = pc.Delete(ctxVW)
	if err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

func (dbc *DBConnector) Select(ctx context.Context, t string) (model.Appender, error) {

	user := ctx.Value(model.KeyContext("user"))
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return nil, errs.ErrErrorServer
	}
	defer conn.Release()

	strUser := user.(string)
	na, err := model.NewAppender(t, strUser)
	if err != nil {
		return nil, errs.ErrErrorServer
	}

	ctxVW := context.WithValue(ctx, model.KeyContext("data"), na)
	pc := model.PgxpoolConn{conn}

	a, err := pc.Select(ctxVW)

	if err != nil {
		return nil, errs.ErrErrorServer
	}

	return a, nil
}

func (dbc *DBConnector) Update(u model.Updater) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	ctxVW := context.WithValue(ctx, model.KeyContext("data"), u)
	pc := model.PgxpoolConn{conn}

	recordExists, err := pc.CheckExistence(ctxVW)

	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {

		if err = pc.Update(ctxVW); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}

	if err = pc.Insert(ctxVW); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) Delete(u model.Updater) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	ctxVW := context.WithValue(ctx, model.KeyContext("data"), u)
	pc := model.PgxpoolConn{conn}

	err = pc.Delete(ctxVW)
	if err != nil {
		return errs.ErrErrorServer
	}
	return nil
}

func (dbc *DBConnector) Eventing(u model.Updater) string {
	return u.GetEvent()
}

/////////////////////////////////////

func (dbc *DBConnector) SelectPortionBinaryData(ctx context.Context) ([]*model.PortionBinaryData, error) {

	uid := ctx.Value(model.KeyContext("uid"))
	rows, err := dbc.Pool.Query(ctx, constants.QuerySelectPortionsBinaryData, uid)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	var arrPbd []*model.PortionBinaryData
	for rows.Next() {
		var pbd *model.PortionBinaryData

		err = rows.Scan(&pbd.Uid, &pbd.Portion, &pbd.Body)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}
		arrPbd = append(arrPbd, pbd)
	}

	return arrPbd, nil
}
