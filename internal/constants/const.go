// Package constants: константы
package constants

import (
	"time"

	"gophkeeper/internal/logger"

	"github.com/gdamore/tcell/v2"
)

type TypeRecord int
type EventDB int

const (
	// TypePairLoginPassword тип хранимой - информации пары логин/пароль
	TypePairLoginPassword TypeRecord = iota

	// TypeTextData тип хранимой информации - произвольные текстовые данные
	TypeTextData

	// TypeBinaryData тип хранимой информации - произвольные бинарные данные
	TypeBinaryData

	// TypeBankCardData тип хранимой информации - данные банковских карт
	TypeBankCardData

	// TypeUserData тип хранимой информации - пользователи
	TypeUserData

	// TypeAuthorizationData тип информации - авторизация пользователя
	TypeAuthorizationData
)

const (
	//EventAddEdit действие с данными добавление/редактирование
	EventAddEdit EventDB = iota

	//EventDel действие с данными удаление
	EventDel
)

const (
	// AdressServer адрес сервера по умолчанию
	AdressServer = "localhost:8080"

	// HeaderAuthorization ключ хедера с авторизированным пользователем
	HeaderAuthorization = "Authorization"

	// Step размер отрезков в байтах, на который "режим" файл
	Step = 512000

	// DefaultColorClient цвет шрифта клиенского приложения
	DefaultColorClient = tcell.ColorGreen

	//NameMainPage имя основного окна клиентского приложения
	NameMainPage = "Menu"
)

const (
	//QuerySelectUserWithWhereTemplate запрос на выборку пользователя по имени
	QuerySelectUserWithWhereTemplate = `SELECT 
								* 
							FROM 
								gophkeeper."Users"
							WHERE 
								"User" = $1;`

	//QuerySelectUserWithPassword запрос на выборку пользователя по имени и паролю
	QuerySelectUserWithPassword = `SELECT 
								* 
							FROM 
								gophkeeper."Users"
							WHERE 
								"User" = $1 AND "Password" = $2;`

	//QueryInsertUserTemplate запрос на добавление пользователя по имени
	QueryInsertUserTemplate = `INSERT INTO 
								gophkeeper."Users" ("User", "Password") 
							VALUES
								($1, $2);`

	//QueryDeleteUserTemplate удаление пользователя
	QueryDeleteUserTemplate = `DELETE 
							FROM 
								gophkeeper."Users"
							WHERE 
								"User" = $1 and "Password" = $2`

	//QueryUpdatUserTemplate запрос на изменение пользователя по имени
	QueryUpdatUserTemplate = `UPDATE gophkeeper."Users" ("User", "Password")
							SET "User"=$1, "Password"=$2
							WHERE 
								"User" = $1;`
) //User

const (
	//QueryInsertPairsTemplate запрос на добавление пары логин/пароль
	QueryInsertPairsTemplate = `INSERT INTO gophkeeper."PairsLoginPassword"(
								"User", "UID", "TypePairs", "Name", "Password")
							VALUES ($1, $2, $3, $4, $5);`

	//QueryUpdatePairsTemplate запрос на изменение пары логин/пароль по пользователю и УИДу
	QueryUpdatePairsTemplate = `UPDATE gophkeeper."PairsLoginPassword"
							SET "User"=$1, "UID"=$2, "TypePairs"=$3, "Name"=$4, "Password"=$5
							WHERE "User" = $1 and "UID" = $2;`

	//QuerySelectPairsTemplate запрос на выборку пары логин/пароль по пользователю
	QuerySelectPairsTemplate = `SELECT "User", "UID", "TypePairs", "Name", "Password" 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1;`

	//QuerySelectOnePairsTemplate запрос на выборку пары логин/пароль по пользователю и УИДу
	QuerySelectOnePairsTemplate = `SELECT "User", "UID", "TypePairs", "Name", "Password" 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1 and "UID" = $2;`

	//QueryDelOnePairsTemplate запрос на уделению пары логин/пароль по пользователю и УИДу
	QueryDelOnePairsTemplate = `DELETE 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1 and "UID" = $2;`
) //Pairs

const (
	//QueryInsertTextData запрос на добавление произвольных текстовых данных
	QueryInsertTextData = `INSERT INTO gophkeeper."Text"(
								"User", "UID", "Text")
							VALUES ($1, $2, $3);`

	//QueryUpdateTextData запрос на изменение произвольных текстовых данных по пользователю и УИДу
	QueryUpdateTextData = `UPDATE gophkeeper."Text"
								SET "User"=$1, "UID"=$2, "Text"=$3 
								WHERE "User" = $1 and "UID" = $2;`

	//QuerySelectTextData запрос на выборку произвольных текстовых данных по пользователю
	QuerySelectTextData = `SELECT "User", "UID", "Text" 
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1;`

	//QuerySelectOneTextData запрос на выборку произвольных текстовых данных по пользователю и УИДу
	QuerySelectOneTextData = `SELECT "User", "UID", "Text" 	
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1 and "UID" = $2;`

	//QueryDelOneTextDataTemplate запрос на уделению произвольных текстовых данных по пользователю и УИДу
	QueryDelOneTextDataTemplate = `DELETE 
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //TextData

const (
	//QueryInsertBankCard запрос на добавление данных банковских карт
	QueryInsertBankCard = `INSERT INTO gophkeeper."BankCards"(
								"User", "UID", "Number", "Cvc")
							VALUES ($1, $2, $3, $4);`

	//QueryUpdateBankCard запрос на изменение данных банковских карт по пользователю и УИДу
	QueryUpdateBankCard = `UPDATE gophkeeper."BankCards"
								SET "User"=$1, "UID"=$2, "Number"=$3, "Cvc"=$4 
								WHERE "User" = $1 and "UID" = $2;`

	//QuerySelectBankCard запрос на выборку данных банковских карт по пользователю
	QuerySelectBankCard = `SELECT "User", "UID", "Number", "Cvc" 
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1;`

	//QuerySelectOneBankCard запрос на выборку данных банковских карт по пользователю и УИДу
	QuerySelectOneBankCard = `SELECT "User", "UID", "Number", "Cvc" 	
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1 and "UID" = $2;`

	//QueryDelOneBankCardTemplate запрос на уделению данных банковских карт по пользователю и УИДу
	QueryDelOneBankCardTemplate = `DELETE 
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //BankCard

const (
	//QueryInsertBinaryData запрос на добавление произвольных бинарных данных
	QueryInsertBinaryData = `INSERT INTO gophkeeper."Files"(
								"User", "UID", "Name", "Expansion", "Size", "Patch")
							VALUES ($1, $2, $3, $4, $5, $6);`

	//QueryUpdateBinaryData запрос на изменение произвольных бинарных данных по пользователю и УИДу
	QueryUpdateBinaryData = `UPDATE gophkeeper."Files"
								SET "User" = $1, "UID" = $2, "Name" = $3, "Expansion" = $4, "Size" = $5, "Patch" = $6 
								WHERE "User" = $1 and "UID" = $2;`

	//QuerySelectBinaryData запрос на выборку произвольных бинарных данных по пользователю
	QuerySelectBinaryData = `SELECT "User", "UID", "Name", "Expansion", "Size", "Patch" 
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1;`

	//QuerySelectOneBinaryData запрос на выборку произвольных бинарных данных по пользователю и УИДу
	QuerySelectOneBinaryData = `SELECT "User", "UID", "Name", "Expansion", "Size", "Patch" 	
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1 and "UID" = $2;`

	//QueryDelOneBinaryDataTemplate запрос на уделению произвольных бинарных данных по пользователю и УИДу
	QueryDelOneBinaryDataTemplate = `DELETE 
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //BinaryData

const (
	//QuerySelectPortionsBinaryData запрос на выборку файлов для таблицы бинарных данных по УИДу
	QuerySelectPortionsBinaryData = `SELECT
							"UID", "Portion", "Body"
						FROM
							gophkeeper."PortionsFiles"
						WHERE
							"UID" = $1;`

	//QueryInsertPortionsBinaryData запрос на добавление файлов для таблицы бинарных данных
	QueryInsertPortionsBinaryData = `INSERT INTO 
							gophkeeper."PortionsFiles"("UID", "Portion", "Body")
						VALUES ($1, $2, $3);`

	//QueryDelPortionsBinaryData запрос на уделению файлов для таблицы бинарных данных по УИДу
	QueryDelPortionsBinaryData = `DELETE FROM gophkeeper."PortionsFiles"	
						WHERE 
							"UID" = $1;`
) //PortionsBinaryData

const (
	KeyCtrlC = 3
	Key0     = 48
	Key1     = 49
	Key2     = 50
	Key3     = 51
	Key4     = 52
	Key5     = 53
	Key6     = 54
	Key7     = 55
)

// HashKey ключ по умолчанию для хешированию паролей
var HashKey = []byte("taekwondo")

// TimeLiveToken время жизни токена. После завершения времени надо перелогиниться.
var TimeLiveToken time.Duration = 5

// Logger логер системы
var Logger logger.Logger

// String  func (tr TypeRecord) String() string преобразует тип хранимой информации в строку
func (tr TypeRecord) String() string {
	return [...]string{"Pairs login/password", "Text", "Binary", "Bank card", "Users", "User authorization"}[tr]
}

// String  func (e EventDB) String() string string преобразует действие с информацией в строку
func (e EventDB) String() string {
	return [...]string{"edit", "del"}[e]
}
