package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/token"
)

// TextData объект текстовые данные
type TextData struct {
	User  string `json:"user"`
	Uid   string `json:"uid"`
	Text  string `json:"text"`
	Event string `json:"event"`
}

// CheckExistence метод объекта TextData. Возвращает инструкции для проверки на существование в БД,
// по пользователю и УИДу
func (t *TextData) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid}
	return constants.QuerySelectOneTextData, arg, nil
}

// InstructionsSelect метод объекта TextData. Возвращает инструкции для добавления баковсокй карты в БД
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

// InstructionsInsert метод объекта TextData. Возвращает инструкции для добавления объекта в БД
func (t *TextData) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid, t.Text}
	return constants.QueryInsertTextData, arg, nil
}

// InstructionsUpdate метод объекта TextData. Возвращает инструкции для обновления объекта в БД,
// по пользователю и УИДу
func (t *TextData) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(t.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], t.Uid, t.Text}
	return constants.QueryUpdateTextData, arg, nil
}

// InstructionsDelete метод объекта TextData. Возвращает инструкции для удаления объекта из БД,
// по пользователю и УИДу
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

// GetEvent метод объекта TextData. Возвращает событие, которое должно произойти с объектом
// в БД. Удаление или добавление/обновление
func (t *TextData) GetEvent() string {
	return t.Event
}

// SetValue метод добавляет объект TextData во временное хранилище сервера
func (t *TextData) SetValue(a Appender) {
	a[t.Uid] = t
}
