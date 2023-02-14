// Package model: работа с моделями базы данных
package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/token"
	"time"
)

// BankCard объект банковская карта
type BankCard struct {
	User   string        `json:"user"`
	Uid    string        `json:"uid"`
	Number string        `json:"patch"`
	Date   time.Duration `json:"date,omitempty"`
	Cvc    string        `json:"cvc"`
	Event  string        `json:"event"`
}

// CheckExistence метод объекта BankCard. Возвращает инструкции для проверки на существование в БД,
// по пользователю и УИДу
func (b *BankCard) CheckExistence() (string, interface{}, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid}
	return constants.QuerySelectOneBankCard, arg, nil
}

// InstructionsSelect метод объекта BankCard. Возвращает инструкции для добавления баковсокй карты в БД
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

// InstructionsInsert метод объекта BankCard. Возвращает инструкции для добавления объекта в БД
func (b *BankCard) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Number, b.Cvc}
	return constants.QueryInsertBankCard, arg, nil
}

// InstructionsUpdate метод объекта BankCard. Возвращает инструкции для обновления объекта в БД, по пользователю и УИДу
func (b *BankCard) InstructionsUpdate() (string, interface{}, error) {

	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Number, b.Cvc}
	return constants.QueryUpdateBankCard, arg, nil
}

// InstructionsDelete метод объекта BankCard. Возвращает инструкции для удаления объекта из БД, по пользователю и УИДу
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

// GetEvent метод объекта BankCard. Возвращает событие, которое должно произойти с объектом
// в БД. Удаление или добавление/обновление
func (b *BankCard) GetEvent() string {
	return b.Event
}

// SetValue метод добавляет объект BankCard во временное хранилище сервера
func (b *BankCard) SetValue(a Appender) {
	a[b.Uid] = b
}
