package model

import (
	"context"
	"errors"
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"

	"github.com/jackc/pgx/v4/pgxpool"
)

type KeyContext string

type Updater interface {
	CheckExistence() (string, interface{}, error)
	InstructionsUpdate() (string, interface{}, error)
	InstructionsInsert() (string, interface{}, error)
	InstructionsDelete() ([]ActionDatabase, error)
	InstructionsSelect() (ActionDatabase, error)

	ReaderWriter
}

type ReaderWriter interface {
	IReader
	IWriter
}

type IReader interface {
	SetValue(Appender)
}

type IWriter interface {
	GetEvent() string
	GetType() string
	GetMainText() string
	GetSecondaryText(string) string
}

type PgxpoolConn struct {
	*pgxpool.Conn
}

type Appender map[string]Updater

type PortionBinaryData struct {
	Type    string `json:"type"`
	Uid     string `json:"uid"`
	Portion int64  `json:"portion"`
	Body    string `json:"body"`
}

type UpenderOut struct {
	Updater
	ArgOut []interface{}
}

func (conn *PgxpoolConn) CheckExistence(ctx context.Context) (bool, error) {

	valCtx := ctx.Value(KeyContext("data"))
	b := valCtx.(Updater)
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

func (conn *PgxpoolConn) Insert(ctx context.Context) error {

	valCtx := ctx.Value(KeyContext("data"))
	b := valCtx.(Updater)
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

func (conn *PgxpoolConn) Update(ctx context.Context) error {

	valCtx := ctx.Value(KeyContext("data"))
	b := valCtx.(Updater)
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

func (conn *PgxpoolConn) Delete(ctx context.Context) error {

	valCtx := ctx.Value(KeyContext("data"))
	b := valCtx.(Updater)
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

func (conn *PgxpoolConn) Select(ctx context.Context) (Appender, error) {

	valCtx := ctx.Value(KeyContext("data"))
	v := valCtx.(Updater)
	actionDatabase, err := v.InstructionsSelect()
	if err != nil {
		return nil, err
	}

	arg := actionDatabase.Arg
	rows, err := conn.Query(ctx, actionDatabase.StrExec, arg...)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	var appender = Appender{}

	for rows.Next() {
		app, err := NewAppender(actionDatabase.Type, actionDatabase.User)
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

func NewAppender(t, u string) (UpenderOut, error) {
	switch t {
	case constants.TypePairLoginPassword.String():
		p := &PairLoginPassword{User: u}
		return UpenderOut{p, []interface{}{&p.User, &p.Uid, &p.TypePair, &p.Name, &p.Password}}, nil
	case constants.TypeTextData.String():
		t := &TextData{User: u}
		return UpenderOut{t, []interface{}{&t.User, &t.Uid, &t.Text}}, nil
	case constants.TypeBinaryData.String():
		b := &BinaryData{User: u}
		return UpenderOut{b, []interface{}{&b.User, &b.Uid, &b.Name, &b.Expansion, &b.Size, &b.Patch}}, nil
	case constants.TypeBankCardData.String():
		b := &BankCard{User: u}
		return UpenderOut{b, []interface{}{&b.User, &b.Uid, &b.Number, &b.Cvc}}, nil
	case constants.TypeUserData.String():
		u := &User{Name: u}
		return UpenderOut{u, []interface{}{&u.Name, &u.Password}}, nil
	default:
		return UpenderOut{}, errors.New("ошибка определения типа данных")
	}
}
