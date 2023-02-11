package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/token"
	"time"
)

type PairLoginPassword struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Uid      string `json:"uid"`
	TypePair string `json:"type_pair"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Event    string `json:"event"`
}

type TextData struct {
	Type  string `json:"type"`
	User  string `json:"user"`
	Uid   string `json:"uid"`
	Text  string `json:"text"`
	Event string `json:"event"`
}

type BinaryData struct {
	Type          string `json:"type"`
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
	Type   string        `json:"type"`
	User   string        `json:"user"`
	Uid    string        `json:"uid"`
	Number string        `json:"patch"`
	Date   time.Duration `json:"date,omitempty"`
	Cvc    string        `json:"cvc"`
	Event  string        `json:"event"`
}

func (p *PairLoginPassword) SetValue(a Appender) {
	a[p.Uid] = p
}

func (p *PairLoginPassword) InstructionsSelect() (ActionDatabase, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return ActionDatabase{}, errs.ErrInvalidLoginPassword
	}

	actionDatabase := ActionDatabase{
		StrExec: constants.QuerySelectPairsTemplate,
		Arg:     []interface{}{claims["user"]},
		Type:    p.GetType(),
		User:    claims["user"].(string),
	}

	//[]interface{}{&p.User, &p.Uid, &p.TypePair, &p.Name, &p.Password},
	return actionDatabase, nil
}

// CheckExistence метод объекта PairLoginPassword проверяющий на существование в БД, по пользователю и УИДу
func (p *PairLoginPassword) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid}
	return constants.QuerySelectOnePairsTemplate, arg, nil
}

// InstructionsInsert метод объекта PairLoginPassword. Добавляет объект в БД
func (p *PairLoginPassword) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid, p.TypePair, p.Name, p.Password}
	return constants.QueryInsertPairsTemplate, arg, nil
}

// InstructionsUpdate метод объекта PairLoginPassword. Обновляет объект в БД, по пользователю и УИДу
func (p *PairLoginPassword) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid, p.TypePair, p.Name, p.Password}
	return constants.QueryUpdatePairsTemplate, arg, nil
}

// InstructionsDelete метод объекта PairLoginPassword. Удаляет объект в БД, по пользователю и УИДу
func (p *PairLoginPassword) InstructionsDelete() ([]ActionDatabase, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return nil, errs.ErrInvalidLoginPassword
	}

	arrActionDatabase := []ActionDatabase{}
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDelOnePairsTemplate,
		Arg:     []interface{}{claims["user"], p.Uid},
	})

	return arrActionDatabase, nil
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
	return encryption.DecryptString(p.TypePair, cryptoKey) + ":::" +
		encryption.DecryptString(p.Name, cryptoKey) + ":::" +
		encryption.DecryptString(p.Password, cryptoKey)
}

// SetFromInListUserData метод объекта PairLoginPassword. Добавляет оьъект в хранилище сервера InListUserData
func (p *PairLoginPassword) SetFromInListUserData(a Appender) {
	a[p.Uid] = p
}

func (p *PairLoginPassword) GetEvent() string {
	return p.Event
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *TextData) SetValue(a Appender) {
	a[t.Uid] = t
}

func (t *TextData) InstructionsSelect() (ActionDatabase, error) {

	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return ActionDatabase{}, errs.ErrInvalidLoginPassword
	}

	actionDatabase := ActionDatabase{
		StrExec: constants.QuerySelectTextData,
		Arg:     []interface{}{claims["user"]},
		Type:    t.GetType(),
		User:    claims["user"].(string),
	}

	//[]interface{}{&t.User, &t.Uid, &t.Text},
	return actionDatabase, nil
}

// InstructionsUpdate метод объекта TextData. Обновляет объект в БД, по пользователю и УИДу
func (t *TextData) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid, t.Text}
	return constants.QueryUpdateTextData, arg, nil
}

// CheckExistence метод объекта TextData проверяющий на существование в БД, по пользователю и УИДу
func (t *TextData) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid}
	return constants.QuerySelectOneTextData, arg, nil
}

// InstructionsInsert метод объекта TextData. Добавляет объект в БД
func (t *TextData) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid, t.Text}
	return constants.QueryInsertTextData, arg, nil
}

// InstructionsDelete метод объекта TextData. Удаляет объект в БД, по пользователю и УИДу
func (t *TextData) InstructionsDelete() ([]ActionDatabase, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return nil, errs.ErrInvalidLoginPassword
	}

	arrActionDatabase := []ActionDatabase{}
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDelOneTextDataTemplate,
		Arg:     []interface{}{claims["user"], t.Uid},
	})

	return arrActionDatabase, nil
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

func (t *TextData) GetEvent() string {
	return t.Event
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BinaryData) SetValue(a Appender) {
	a[b.Uid] = b
}

func (b *BinaryData) InstructionsSelect() (ActionDatabase, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return ActionDatabase{}, errs.ErrInvalidLoginPassword
	}

	actionDatabase := ActionDatabase{
		StrExec: constants.QuerySelectBinaryData,
		Arg:     []interface{}{claims["user"]},
		Type:    b.GetType(),
		User:    claims["user"].(string),
	}

	//[]interface{}{&b.User, &b.Uid, &b.Name, &b.Expansion, &b.Size, &b.Patch},
	return actionDatabase, nil
}

// CheckExistence метод объекта BinaryData проверяющий на существование в БД, по пользователю и УИДу
func (b *BinaryData) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid}
	return constants.QuerySelectOneBinaryData, arg, nil
}

// InstructionsInsert метод объекта BinaryData. Добавляет объект в БД
func (b *BinaryData) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch}
	return constants.QueryInsertBinaryData, arg, nil
}

// InstructionsUpdate метод объекта BinaryData. Обновляет объект в БД, по пользователю и УИДу
func (b *BinaryData) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch}
	return constants.QueryUpdateBinaryData, arg, nil
}

type ActionDatabase struct {
	StrExec string
	Arg     []interface{}
	Type    string
	User    string
}

// InstructionsDelete метод объекта BinaryData. Удаляет объект в БД, по пользователю и УИДу
func (b *BinaryData) InstructionsDelete() ([]ActionDatabase, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return nil, errs.ErrInvalidLoginPassword
	}

	arrActionDatabase := []ActionDatabase{}
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDelOneBinaryDataTemplate,
		Arg:     []interface{}{claims["user"], b.Uid},
	})
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDelPortionsBinaryData,
		Arg:     []interface{}{b.Uid},
	})

	return arrActionDatabase, nil
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

func (b *BinaryData) GetEvent() string {
	return b.Event
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *BankCard) SetValue(a Appender) {
	a[b.Uid] = b
}

func (b *BankCard) InstructionsSelect() (ActionDatabase, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return ActionDatabase{}, errs.ErrInvalidLoginPassword
	}

	actionDatabase := ActionDatabase{
		StrExec: constants.QuerySelectBankCard,
		Arg:     []interface{}{claims["user"]},
		Type:    b.GetType(),
		User:    claims["user"].(string),
	}

	//[]interface{}{&b.User, &b.Uid, &b.Number, &b.Cvc},
	return actionDatabase, nil
}

// CheckExistence метод объекта BankCard проверяющий на существование в БД, по пользователю и УИДу
func (b *BankCard) CheckExistence() (string, interface{}, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid}
	return constants.QuerySelectOneBankCard, arg, nil
}

// InstructionsInsert метод объекта BankCard. Добавляет объект в БД
func (b *BankCard) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Number, b.Cvc}
	return constants.QueryInsertBankCard, arg, nil
}

// InstructionsUpdate метод объекта BankCard. Обновляет объект в БД, по пользователю и УИДу
func (b *BankCard) InstructionsUpdate() (string, interface{}, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Number, b.Cvc}
	return constants.QueryUpdateBankCard, arg, nil
}

// InstructionsDelete метод объекта BankCard. Удаляет объект в БД, по пользователю и УИДу
func (b *BankCard) InstructionsDelete() ([]ActionDatabase, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return nil, errs.ErrInvalidLoginPassword
	}

	arrActionDatabase := []ActionDatabase{}
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDelOneBankCardTemplate,
		Arg:     []interface{}{claims["user"], b.Uid},
	})

	return arrActionDatabase, nil
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

func (b *BankCard) GetEvent() string {
	return b.Event
}
