package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/token"
)

type User struct {
	Type         string `json:"type"`
	Name         string `json:"login"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password"`
	Event        string `json:"event"`
}

// GetType метод объекта PairLoginPassword. Возвращает текстовое представление объекта
func (u *User) GetType() string {
	return constants.TypePairLoginPassword.String()
}

// GetMainText метод объекта TextData. Создает основной текст для объекта List, клиентского приложения.
func (u *User) GetMainText() string {
	return ""
}

// GetSecondaryText метод объекта TextData. Создает вспомогательный текст для объекта List, клиентского приложения.
func (u *User) GetSecondaryText(cryptoKey string) string {
	return encryption.DecryptString("", cryptoKey)
}

func (u *User) SetValue(a Appender) {
	a[u.Name] = u
}

func (u *User) InstructionsSelect() (ActionDatabase, error) {

	actionDatabase := ActionDatabase{
		StrExec: constants.QuerySelectUserWithWhereTemplate,
		Arg:     []interface{}{u.Name},
		Type:    u.GetType(),
		User:    u.Name,
	}

	return actionDatabase, nil
}

// InstructionsDelete метод объекта User. Удаляет объект из БД по имени и хешированному паролю
func (u *User) InstructionsDelete() ([]ActionDatabase, error) {

	arrActionDatabase := []ActionDatabase{}
	arrActionDatabase = append(arrActionDatabase, ActionDatabase{
		StrExec: constants.QueryDeleteUserTemplate,
		Arg:     []interface{}{u.Name, u.HashPassword},
	})

	return arrActionDatabase, nil
}

// InstructionsInsert метод объекта User. Добавляет объект в БД
func (u *User) InstructionsInsert() (string, interface{}, error) {
	arg := []interface{}{u.Name, u.HashPassword}
	return constants.QueryInsertUserTemplate, arg, nil
}

// CheckExistence метод объекта User проверяющий на существование в БД, по пользователю и УИДу
func (u *User) CheckExistence() (string, interface{}, error) {
	arg := []interface{}{u.Name, u.HashPassword}
	return constants.QuerySelectUserWithPassword, arg, nil
}

// InstructionsUpdate метод объекта PairLoginPassword. Обновляет объект в БД, по пользователю и УИДу
func (u *User) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(u.Name)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], u.Password}
	return constants.QueryUpdatUserTemplate, arg, nil
}

// SetFromInListUserData метод объекта User. Добавляет оьъект в хранилище сервера InListUserData
func (u *User) SetFromInListUserData(a Appender) {
	a[u.Name] = u
}

func (u *User) GetEvent() string {
	return u.Event
}
