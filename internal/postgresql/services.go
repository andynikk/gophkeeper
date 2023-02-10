package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/cryptography"
)

type Updater interface {
	CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error)
	Update(ctx context.Context, conn *pgxpool.Conn) error
	Insert(ctx context.Context, conn *pgxpool.Conn) error
	Delete(ctx context.Context, conn *pgxpool.Conn) error
	Select(ctx context.Context, conn *pgxpool.Conn) (Appender, error)

	SetValue(Appender)
	GetEvent() string

	GetType() string
	GetMainText() string
	GetSecondaryText(string) string
}

// NewAccount метод для создания нового экаунта из ДБ конектора
// вызывает методы объекта user.
// Проверяет есть ли такой пользователь.
// Если нет, то создает
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

// CheckAccount проверяет, по имени и хешированному паролю существует ли пользователь в базе
func (dbc *DBConnector) CheckAccount(user User) error {

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

// DelAccount удаляет пользователя по имени и хешированному паролю
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

func NewAppender(t, u string) (Updater, error) {
	switch t {
	case constants.TypePairLoginPassword.String():
		return &PairLoginPassword{User: u}, nil
	case constants.TypeTextData.String():
		return &TextData{User: u}, nil
	case constants.TypeBinaryData.String():
		return &BinaryData{User: u}, nil
	case constants.TypeBankCardData.String():
		return &BankCard{User: u}, nil
	case constants.TypeUserData.String():
		return &User{Name: u}, nil
	default:
		return nil, errors.New("ошибка определения типа данных")
	}
}

func (dbc *DBConnector) Select(ctx context.Context, t string) (Appender, error) {

	user := ctx.Value(KeyContext("user"))
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return nil, errs.ErrErrorServer
	}
	defer conn.Release()

	strUser := user.(string)
	na, err := NewAppender(t, strUser)
	if err != nil {
		return nil, errs.ErrErrorServer
	}
	a, err := na.Select(ctx, conn)
	if err != nil {
		return nil, errs.ErrErrorServer
	}

	return a, nil
}

func (dbc *DBConnector) Update(u Updater) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	recordExists, err := u.CheckExistence(ctx, conn)
	if err != nil {
		return errs.InvalidFormat
	}
	if recordExists {
		if err = u.Update(ctx, conn); err != nil {
			return errs.InvalidFormat
		}
		return nil
	}

	if err = u.Insert(ctx, conn); err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (dbc *DBConnector) Delete(u Updater) error {
	ctx := context.Background()
	conn, err := dbc.Pool.Acquire(ctx)
	if err != nil {
		return errs.ErrErrorServer
	}
	defer conn.Release()

	err = u.Delete(ctx, conn)
	if err != nil {
		return errs.ErrErrorServer
	}
	return nil
}

func (dbc *DBConnector) Eventing(u Updater) string {
	return u.GetEvent()
}

/////////////////////////////////////

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
