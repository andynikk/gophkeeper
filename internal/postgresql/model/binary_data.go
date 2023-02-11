package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/token"
)

// BinaryData объект бинарные данные
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

// CheckExistence метод объекта BinaryData. Возвращает инструкции для проверки на существование в БД,
// по пользователю и УИДу
func (b *BinaryData) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid}
	return constants.QuerySelectOneBinaryData, arg, nil
}

// InstructionsSelect метод объекта BinaryData. Возвращает инструкции для добавления баковсокй карты в БД
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

	return actionDatabase, nil
}

// InstructionsInsert метод объекта BinaryData. Возвращает инструкции для добавления объекта в БД
func (b *BinaryData) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch}
	return constants.QueryInsertBinaryData, arg, nil
}

// InstructionsUpdate метод объекта BinaryData. Возвращает инструкции для обновления объекта в БД,
// по пользователю и УИДу
func (b *BinaryData) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(b.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], b.Uid, b.Name, b.Expansion, b.Size, b.Patch}
	return constants.QueryUpdateBinaryData, arg, nil
}

// InstructionsDelete метод объекта BinaryData. Возвращает инструкции для удаления объекта из БД,
// по пользователю и УИДу
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

// GetEvent метод объекта BinaryData. Возвращает событие, которое должно произойти с объектом
// в БД. Удаление или добавление/обновление
func (b *BinaryData) GetEvent() string {
	return b.Event
}

// SetFromInListUserData метод объекта BinaryData. Добавляет оьъект в хранилище сервера InListUserData
func (b *BinaryData) SetFromInListUserData(a Appender) {
	a[b.Uid] = b
}

// SetValue метод добавляет объект BinaryData во временное хранилище сервера
func (b *BinaryData) SetValue(a Appender) {
	a[b.Uid] = b
}
