package postgresql

import (
	"context"
	"gophkeeper/internal/token"
	"time"

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

type PairLoginPassword struct {
	User      string `json:"user"`
	Uid       string `json:"uid"`
	TypePairs string `json:"type"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Event     string `json:"event"`
}

type TextData struct {
	User  string `json:"user"`
	Uid   string `json:"uid"`
	Text  string `json:"text"`
	Event string `json:"event"`
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
	Event         string `json:"event"`
}

type BankCard struct {
	User   string        `json:"user"`
	Uid    string        `json:"uid"`
	Number string        `json:"patch"`
	Date   time.Duration `json:"date,omitempty"`
	Cvc    string        `json:"cvc"`
	Event  string        `json:"event"`
}

type User struct {
	Name         string `json:"login"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password"`
	Event        string `json:"event"`
}

type Users struct {
	Users []User
}

type UsersWithType struct {
	Type  string
	Event string
	Value []User
}

type PairLoginPasswordWithType struct {
	Type  string
	Event string
	Value []PairLoginPassword
}

type TextDataWithType struct {
	Type     string
	Event    string
	TextData []TextData
}

type BinaryDataWithType struct {
	Type       string
	Event      string
	BinaryData []BinaryData
}

type BankCardWithType struct {
	Type     string
	Event    string
	BankCard []BankCard
}

type Appender map[string]interface{}

type TypeMsg struct {
	Type  string
	Token string
}

type InListUserData interface {
	SetFromInListUserData(Appender)
}

type DataList struct {
	TypeResponse  string `json:"type"`
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}

type KeyContext string

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckExistence метод объекта PairLoginPassword проверяющий на существование в БД, по пользователю и УИДу
func (p *PairLoginPassword) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return false, errs.ErrInvalidLoginPassword
	}

	rows, err := conn.Query(ctx, constants.QuerySelectOnePairsTemplate, claims["user"], p.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Insert метод объекта PairLoginPassword. Добавляет объект в БД
func (p *PairLoginPassword) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryInsertPairsTemplate, claims["user"], p.Uid, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Update метод объекта PairLoginPassword. Обновляет объект в БД, по пользователю и УИДу
func (p *PairLoginPassword) Update(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryUpdatePairsTemplate, claims["user"], p.Uid, p.TypePairs, p.Name, p.Password)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// GetType метод объекта PairLoginPassword. Возвращает текстовое представление объекта
func (p *PairLoginPassword) GetType() string {
	return constants.TypePairLoginPassword.String()
}

// GetMainText метод объекта PairLoginPassword. Создает основной текст для объекта List, клиентского приложения.
func (p *PairLoginPassword) GetMainText() string {
	return p.Uid
}

// GetSecondaryText метод объекта PairLoginPassword. Создает вспомогательный текст для объекта List, клиентского приложения.
func (p *PairLoginPassword) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(p.TypePairs, cryptoKey) + ":::" +
		encryption.DecryptString(p.Name, cryptoKey) + ":::" +
		encryption.DecryptString(p.Password, cryptoKey)
}

// SetFromInListUserData метод объекта PairLoginPassword. Добавляет оьъект в хранилище сервера InListUserData
func (p *PairLoginPassword) SetFromInListUserData(a Appender) {
	a[p.Uid] = p
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Delete метод объекта User. Удаляет объект из БД по имени и хешированному паролю
func (u *User) Delete(ctx context.Context, conn *pgxpool.Conn) error {
	if _, err := conn.Exec(ctx, constants.QueryDeleteUserTemplate, u.Name, u.HashPassword); err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

// Insert метод объекта User. Добавляет объект в БД
func (u *User) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	if _, err := conn.Exec(ctx, constants.QueryInsertUserTemplate, u.Name, u.HashPassword); err != nil {
		return errs.ErrErrorServer
	}

	return nil
}

// CheckExistence метод объекта User проверяющий на существование в БД, по пользователю и УИДу
func (u *User) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	rows, err := conn.Query(ctx, constants.QuerySelectUserWithWhereTemplate, u.Name)
	if err != nil {
		return false, errs.ErrErrorServer
	}
	defer rows.Close()

	return rows.Next(), nil
}

// SetFromInListUserData метод объекта User. Добавляет оьъект в хранилище сервера InListUserData
func (u *User) SetFromInListUserData(a Appender) {
	a[u.Name] = u
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckExistence метод объекта TextData проверяющий на существование в БД, по пользователю и УИДу
func (t *TextData) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return false, errs.ErrInvalidLoginPassword
	}

	rows, err := conn.Query(ctx, constants.QuerySelectOneTextData, claims["user"], t.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Insert метод объекта TextData. Добавляет объект в БД
func (t *TextData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryInsertTextData, claims["user"], t.Uid, t.Text)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Update метод объекта TextData. Обновляет объект в БД, по пользователю и УИДу
func (t *TextData) Update(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryUpdateTextData, claims["user"], t.Uid, t.Text)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// GetType метод объекта TextData. Возвращает текстовое представление объекта
func (t *TextData) GetType() string {
	return constants.TypeTextData.String()
}

// GetMainText метод объекта TextData. Создает основной текст для объекта List, клиентского приложения.
func (t *TextData) GetMainText() string {
	return t.Uid
}

// GetSecondaryText метод объекта TextData. Создает вспомогательный текст для объекта List, клиентского приложения.
func (t *TextData) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(t.Text, cryptoKey)
}

// SetFromInListUserData метод объекта TextData. Добавляет оьъект в хранилище сервера InListUserData
func (t *TextData) SetFromInListUserData(a Appender) {
	a[t.Uid] = t
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckExistence метод объекта BinaryData проверяющий на существование в БД, по пользователю и УИДу
func (b *BinaryData) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return false, errs.ErrInvalidLoginPassword
	}

	rows, err := conn.Query(ctx, constants.QuerySelectOneBinaryData, claims["user"], b.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Insert метод объекта BinaryData. Добавляет объект в БД
func (b *BinaryData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryInsertBinaryData,
		claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Update метод объекта BinaryData. Обновляет объект в БД, по пользователю и УИДу
func (b *BinaryData) Update(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryUpdateBinaryData, claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// GetType метод объекта BinaryData. Возвращает текстовое представление объекта
func (b *BinaryData) GetType() string {
	return constants.TypeBinaryData.String()
}

// GetMainText метод объекта BinaryData. Создает основной текст для объекта List, клиентского приложения.
func (b *BinaryData) GetMainText() string {
	return b.Uid
}

// GetSecondaryText метод объекта BinaryData. Создает вспомогательный текст для объекта List, клиентского приложения.
func (b *BinaryData) GetSecondaryText(cryptoKey string) string {
	return b.Name + ":::" + b.Expansion + ":::" + b.Size + ":::" + b.Patch
}

// SetFromInListUserData метод объекта BinaryData. Добавляет оьъект в хранилище сервера InListUserData
func (b *BinaryData) SetFromInListUserData(a Appender) {
	a[b.Uid] = b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckExistence метод объекта BankCard проверяющий на существование в БД, по пользователю и УИДу
func (b *BankCard) CheckExistence(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return false, errs.ErrInvalidLoginPassword
	}

	rows, err := conn.Query(ctx, constants.QuerySelectOneBankCard, claims["user"], b.Uid)
	if err != nil {
		return false, errs.InvalidFormat
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Insert метод объекта BankCard. Добавляет объект в БД
func (b *BankCard) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryInsertBankCard, claims["user"], b.Uid, b.Number, b.Cvc)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// Update метод объекта BankCard. Обновляет объект в БД, по пользователю и УИДу
func (b *BankCard) Update(ctx context.Context, conn *pgxpool.Conn) error {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return errs.ErrInvalidLoginPassword
	}

	_, err := conn.Exec(ctx, constants.QueryUpdateBankCard, claims["user"], b.Uid, b.Number, b.Cvc)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}

// GetType метод объекта BankCard. Возвращает текстовое представление объекта
func (b *BankCard) GetType() string {
	return constants.TypeBankCardData.String()
}

// GetMainText метод объекта BankCard. Создает основной текст для объекта List, клиентского приложения.
func (b *BankCard) GetMainText() string {
	return b.Uid
}

// GetSecondaryText метод объекта BankCard. Создает вспомогательный текст для объекта List, клиентского приложения.
func (b *BankCard) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString(b.Number, cryptoKey) + ":::" +
		encryption.DecryptString(b.Cvc, cryptoKey)
}

// SetFromInListUserData метод объекта BankCard. Добавляет оьъект в хранилище сервера InListUserData
func (b *BankCard) SetFromInListUserData(a Appender) {
	a[b.Uid] = b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Insert метод объекта PortionBinaryData. Добавляет объект в БД
func (p *PortionBinaryData) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	_, err := conn.Exec(ctx, constants.QueryInsertPortionsBinaryData, p.Uid, p.Portion, p.Body)
	if err != nil {
		return errs.InvalidFormat
	}

	return nil
}
