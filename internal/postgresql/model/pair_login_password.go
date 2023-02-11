package model

import (
	"gophkeeper/internal/constants"
	"gophkeeper/internal/constants/errs"
	"gophkeeper/internal/encryption"
	"gophkeeper/internal/token"
)

// PairLoginPassword объект пара логин/пароль
type PairLoginPassword struct {
	User     string `json:"user"`
	Uid      string `json:"uid"`
	TypePair string `json:"type_pair"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Event    string `json:"event"`
}

// CheckExistence метод объекта PairLoginPassword. Возвращает инструкции для проверки на существование в БД,
// по пользователю и УИДу
func (p *PairLoginPassword) CheckExistence() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid}
	return constants.QuerySelectOnePairsTemplate, arg, nil
}

// InstructionsSelect метод объекта PairLoginPassword. Возвращает инструкции для добавления баковсокй карты в БД
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

// InstructionsInsert метод объекта PairLoginPassword. Возвращает инструкции для добавления объекта в БД
func (p *PairLoginPassword) InstructionsInsert() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid, p.TypePair, p.Name, p.Password}
	return constants.QueryInsertPairsTemplate, arg, nil
}

// InstructionsUpdate метод объекта PairLoginPassword. Возвращает инструкции для обновления объекта в БД,
// по пользователю и УИДу
func (p *PairLoginPassword) InstructionsUpdate() (string, interface{}, error) {
	claims, ok := token.ExtractClaims(p.User)
	if !ok {
		return "", nil, errs.ErrInvalidLoginPassword
	}

	arg := []interface{}{claims["user"], p.Uid, p.TypePair, p.Name, p.Password}
	return constants.QueryUpdatePairsTemplate, arg, nil
}

// InstructionsDelete метод объекта PairLoginPassword. Возвращает инструкции для удаления объекта из БД,
// по пользователю и УИДу
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

// GetEvent метод объекта PairLoginPassword. Возвращает событие, которое должно произойти с объектом
// в БД. Удаление или добавление/обновление
func (p *PairLoginPassword) GetEvent() string {
	return p.Event
}

// SetValue метод добавляет объект BinaryData во временное хранилище сервера
func (p *PairLoginPassword) SetValue(a Appender) {
	a[p.Uid] = p
}
