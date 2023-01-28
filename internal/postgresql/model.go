package postgresql

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
)

type Respondent interface {
	GetType() string
	GetMainText() string
	GetSecondaryText() string
}

type MapResponse = map[uuid.UUID]Respondent

//
//type HandlerJSON interface {
//	Marshal() ([]byte, error)
//	Unmarshal([]byte) error
//}

//type CheckingExistence interface {
//	CheckExistence(context.Context, *pgxpool.Conn) (bool, error)
//}
//
//type Inserter interface {
//	Insert(context.Context, *pgxpool.Conn) error
//}
//
//type Updater interface {
//	Update(context.Context, *pgxpool.Conn) error
//}

//type Selecter interface {
//	Select(context.Context, *pgxpool.Conn) (interface{}, error)
//}

type User struct {
	Name         string `json:"login"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password"`
}

type PairsLoginPassword struct {
	User      string `json:"user"`
	TypePairs string `json:"type"`
	Name      string `json:"name"`
	Password  string `json:"password"`
}

type TextData struct {
	User string `json:"user"`
	Uid  string `json:"uid"`
	Text string `json:"text"`
}

type KeyContext string

func SelectAll(ctx context.Context, conn *pgxpool.Conn) (MapResponse, error) {

	user := ctx.Value(KeyContext("user"))
	mr := MapResponse{}

	rows, err := conn.Query(ctx, constants.QuerySelectPairsTemplate, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}
	defer rows.Close()

	for rows.Next() {
		var plp PairsLoginPassword

		err = rows.Scan(&plp.User, &plp.TypePairs, &plp.Name, &plp.Password)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		mr[uuid.New()] = &plp
	}

	rows, err = conn.Query(ctx, constants.QuerySelectTextData, user)
	if err != nil {
		return nil, errs.InvalidFormat
	}

	for rows.Next() {
		var td TextData

		err = rows.Scan(&td.User, &td.Uid, &td.Text)
		if err != nil {
			constants.Logger.ErrorLog(err)
			continue
		}

		mr[uuid.New()] = &td
	}

	return mr, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *PairsLoginPassword) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectOnePairsTemplate, p.User, p.TypePairs, p.Name)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (p *PairsLoginPassword) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertPairsTemplate, p.User, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (p *PairsLoginPassword) Update(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryUpdatePairsTemplate, p.User, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (p *PairsLoginPassword) GetType() string {
	return constants.TypePairsLoginPassword.String()
}

func (p *PairsLoginPassword) GetMainText() string {
	return p.TypePairs
}

func (p *PairsLoginPassword) GetSecondaryText() string {
	return p.Name + ":::" + p.Password
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (u *User) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	if _, err := conn.Exec(ctx, constants.QueryInsertUserTemplate, u.Name, u.HashPassword); err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

func (u *User) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectUserWithWhereTemplate, u.Name)
	if err != nil {
		return false, errs.ErrErrorServer
	}
	defer rows.Close()

	return rows.Next(), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *TextData) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectOneTextData, t.User, t.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (t *TextData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertTextData, t.User, t.Uid, t.Text)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (t *TextData) Update(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryUpdateTextData, t.User, t.Uid, t.Text)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (t *TextData) GetType() string {
	return constants.TypeTextData.String()
}

func (t *TextData) GetMainText() string {
	return t.Uid
}

func (t *TextData) GetSecondaryText() string {
	return t.Text
}
