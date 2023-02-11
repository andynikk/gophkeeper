package model

import (
	"errors"
	"gophkeeper/internal/constants"
)

// KeyContext тип для создания ключа контекста
type KeyContext string

// PortionBinaryData структура кусочка файла
type PortionBinaryData struct {
	Uid     string `json:"uid"`
	Portion int64  `json:"portion"`
	Body    string `json:"body"`
}

// ActionDatabase структура указывающая, что делать с БД, с параметрами для отбора
type ActionDatabase struct {
	StrExec string
	Arg     []interface{}
	Type    string
	User    string
}

// Updater интерфес для работы с объектами БД, отвечающих условиям контракта.
// На текущий момент это User, BankCard, BinaryData, PairLoginPassword, TextData
type Updater interface {
	CheckExistence() (string, interface{}, error)
	InstructionsUpdate() (string, interface{}, error)
	InstructionsInsert() (string, interface{}, error)
	InstructionsDelete() ([]ActionDatabase, error)
	InstructionsSelect() (ActionDatabase, error)

	ReaderWriter
}

// ReaderWriter интерфейс вложение. Использует объекты для чтени и записи
type ReaderWriter interface {
	IReader
	IWriter
}

// IReader интерфейс вложение. Использует объекты для чтения
type IReader interface {
	GetEvent() string
	GetType() string
	GetMainText() string
	GetSecondaryText(string) string
}

// IWriter интерфейс вложение. Использует объекты для записи
type IWriter interface {
	SetValue(Appender)
}

// Appender мапа для хранения объектов Updater
type Appender map[string]Updater

// UpdaterOut структура для создания объекта Updater.
type UpdaterOut struct {
	Updater
	ArgOut []interface{}
}

// NewAppender создание нового объекта Updater.
// С заполненными пользователем, и пустми свойствами для выполнения запроса select
func NewAppender(t, u string) (UpdaterOut, error) {
	switch t {
	case constants.TypePairLoginPassword.String():
		p := &PairLoginPassword{User: u}
		return UpdaterOut{p, []interface{}{&p.User, &p.Uid, &p.TypePair, &p.Name, &p.Password}}, nil
	case constants.TypeTextData.String():
		t := &TextData{User: u}
		return UpdaterOut{t, []interface{}{&t.User, &t.Uid, &t.Text}}, nil
	case constants.TypeBinaryData.String():
		b := &BinaryData{User: u}
		return UpdaterOut{b, []interface{}{&b.User, &b.Uid, &b.Name, &b.Expansion, &b.Size, &b.Patch}}, nil
	case constants.TypeBankCardData.String():
		b := &BankCard{User: u}
		return UpdaterOut{b, []interface{}{&b.User, &b.Uid, &b.Number, &b.Cvc}}, nil
	case constants.TypeUserData.String():
		u := &User{Name: u}
		return UpdaterOut{u, []interface{}{&u.Name, &u.Password}}, nil
	default:
		return UpdaterOut{}, errors.New("ошибка определения типа данных")
	}
}
