package postgresql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
)

type Respondent interface {
	GetType() string
	GetMainText() string
	GetSecondaryText(string) string
}

type MapResponse = map[uuid.UUID]Respondent

type User struct {
	Name         string `json:"login"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password"`
}

type PairsLoginPassword struct {
	User      string `json:"user"`
	Uid       string `json:"uid"`
	TypePairs string `json:"type"`
	Name      string `json:"name"`
	Password  string `json:"password"`
}

type TextData struct {
	User string `json:"user"`
	Uid  string `json:"uid"`
	Text string `json:"text"`
}

type PortionBinaryData struct {
	Uid     string `json:"uid"`
	Portion int64  `json:"portion"`
	Body    string `json:"body"`
}

type BinaryData struct {
	User          string `json:"user"`
	Uid           string `json:"uid"`
	Patch         string `json:"patch"`
	DownloadPatch string `json:"download_patch,omitempty"`
	Name          string `json:"name"`
	Expansion     string `json:"expansion"`
	Size          string `json:"size"`
}

type BankCard struct {
	User   string        `json:"user"`
	Uid    string        `json:"uid"`
	Number string        `json:"patch"`
	Date   time.Duration `json:"date,omitempty"`
	Cvc    string        `json:"cvc"`
}

type PairsLoginPasswordWithType struct {
	Type               string
	PairsLoginPassword []PairsLoginPassword
}

type TextDataWithType struct {
	Type     string
	TextData []TextData
}

type BinaryDataWithType struct {
	Type       string
	BinaryData []BinaryData
}

type BankCardWithType struct {
	Type     string
	BankCard []BankCard
}

type Response struct {
	TypeResponse  string `json:"type"`
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}

type KeyContext string

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *PairsLoginPassword) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectOnePairsTemplate, p.User, p.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (p *PairsLoginPassword) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertPairsTemplate, p.User, p.Uid, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (p *PairsLoginPassword) Update(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryUpdatePairsTemplate, p.User, p.Uid, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (p *PairsLoginPassword) GetType() string {
	return constants.TypePairsLoginPassword.String()
}

func (p *PairsLoginPassword) GetMainText() string {
	return p.Uid
}

func (p *PairsLoginPassword) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(p.TypePairs, cryptoKey) + ":::" +
		encryption.DecryptString(p.Name, cryptoKey) + ":::" +
		encryption.DecryptString(p.Password, cryptoKey)
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

func (t *TextData) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(t.Text, cryptoKey)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BinaryData) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectOneBinaryData, b.User, b.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (b *BinaryData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertBinaryData,
		b.User, b.Uid, b.Name, b.Expansion, b.Size, b.Patch)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (b *BinaryData) Update(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryUpdateBinaryData, b.User, b.Uid, b.Name, b.Expansion, b.Size, b.Patch)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (b *BinaryData) GetType() string {
	return constants.TypeBinaryData.String()
}

func (b *BinaryData) GetMainText() string {
	return b.Uid
}

func (b *BinaryData) GetSecondaryText(cryptoKey string) string {
	return b.Name + ":::" + b.Expansion + ":::" + b.Size + ":::" + b.Patch
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BankCard) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectOneBankCard, b.User, b.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (b *BankCard) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertBankCard, b.User, b.Uid, b.Number, b.Cvc)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (b *BankCard) Update(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryUpdateBankCard, b.User, b.Uid, b.Number, b.Cvc)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

func (b *BankCard) GetType() string {
	return constants.TypeBankCardData.String()
}

func (b *BankCard) GetMainText() string {
	return b.Uid
}

func (b *BankCard) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(b.Number, cryptoKey) + ":::" +
		encryption.DecryptString(b.Cvc, cryptoKey)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (p *PortionBinaryData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertPortionsBinaryData, p.Uid, p.Portion, p.Body)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}
