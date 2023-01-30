package constants

import (
	"time"

	"gophkeeper/internal/logger"
)

type TypeRecord int
type EventDB int

const (
	TypePairsLoginPassword TypeRecord = iota
	TypeTextData
	TypeBinaryData
	TypeBankCardData
)

const (
	EventAddEdit EventDB = iota
	EventDel
)

const (
	AdressServer = "localhost:8080"

	HeaderAuthorization = "Authorization"

	Step = 512000
)

const (
	QuerySelectUserWithWhereTemplate = `SELECT 
								* 
							FROM 
								gophkeeper."Users"
							WHERE 
								"User" = $1;`

	QueryInsertUserTemplate = `INSERT INTO 
								gophkeeper."Users" ("User", "Password") 
							VALUES
								($1, $2);`
) //User

const (
	QueryInsertPairsTemplate = `INSERT INTO gophkeeper."PairsLoginPassword"(
								"User", "UID", "TypePairs", "Name", "Password")
							VALUES ($1, $2, $3, $4, $5);`

	QueryUpdatePairsTemplate = `UPDATE gophkeeper."PairsLoginPassword"
							SET "User"=$1, "UID"=$2, "TypePairs"=$3, "Name"=$4, "Password"=$5
							WHERE "User" = $1 and "UID" = $2;`

	QuerySelectPairsTemplate = `SELECT "User", "UID", "TypePairs", "Name", "Password" 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1;`

	QuerySelectOnePairsTemplate = `SELECT "User", "UID", "TypePairs", "Name", "Password" 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1 and "UID" = $2;`
	QueryDelOnePairsTemplate = `DELETE 
							FROM 
								gophkeeper."PairsLoginPassword"
							WHERE 
								"User" = $1 and ", "UID" = $2;`
) //Pairs

const (
	QueryInsertTextData = `INSERT INTO gophkeeper."Text"(
								"User", "UID", "Text")
							VALUES ($1, $2, $3);`

	QueryUpdateTextData = `UPDATE gophkeeper."Text"
								SET "User"=$1, "UID"=$2, "Text"=$3 
								WHERE "User" = $1 and "UID" = $2;`

	QuerySelectTextData = `SELECT "User", "UID", "Text" 
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1;`

	QuerySelectOneTextData = `SELECT "User", "UID", "Text" 	
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1 and "UID" = $2;`

	QueryDelOneTextDataTemplate = `DELETE 
						FROM 
							gophkeeper."Text"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //TextData

const (
	QueryInsertBankCard = `INSERT INTO gophkeeper."BankCards"(
								"User", "UID", "Number", "Cvc")
							VALUES ($1, $2, $3, $4);`

	QueryUpdateBankCard = `UPDATE gophkeeper."BankCards"
								SET "User"=$1, "UID"=$2, "Number"=$3, "Cvc"=$4 
								WHERE "User" = $1 and "UID" = $2;`

	QuerySelectBankCard = `SELECT "User", "UID", "Number", "Cvc" 
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1;`

	QuerySelectOneBankCard = `SELECT "User", "UID", "Number", "Cvc" 	
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1 and "UID" = $2;`

	QueryDelOneBankCardTemplate = `DELETE 
						FROM 
							gophkeeper."BankCards"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //BankCard

const (
	QueryInsertBinaryData = `INSERT INTO gophkeeper."Files"(
								"User", "UID", "Name", "Expansion", "Size", "Patch")
							VALUES ($1, $2, $3, $4, $5, $6);`

	QueryUpdateBinaryData = `UPDATE gophkeeper."Files"
								SET "User" = $1, "UID" = $2, "Name" = $3, "Expansion" = $4, "Size" = $5, "Patch" = $6 
								WHERE "User" = $1 and "UID" = $2;`

	QuerySelectBinaryData = `SELECT "User", "UID", "Name", "Expansion", "Size", "Patch" 
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1;`

	QuerySelectOneBinaryData = `SELECT "User", "UID", "Name", "Expansion", "Size", "Patch" 	
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1 and "UID" = $2;`

	QueryDelOneBinaryDataTemplate = `DELETE 
						FROM 
							gophkeeper."Files"
						WHERE 
							"User" = $1 and "UID" = $2;`
) //BinaryData

const (
	QuerySelectPortionsBinaryData = `SELECT
							"UID", "Portion", "Body"
						FROM
							gophkeeper."PortionsFiles"
						WHERE
							"UID" = $1;`
	QueryInsertPortionsBinaryData = `INSERT INTO 
							gophkeeper."PortionsFiles"("UID", "Portion", "Body")
						VALUES ($1, $2, $3);`
	QueryDelPortionsBinaryData = `DELETE FROM gophkeeper."PortionsFiles"	
						WHERE 
							"UID" = $1;`
) //PortionsBinaryData

var HashKey = []byte("taekwondo")
var TimeLiveToken time.Duration = 5
var Logger logger.Logger

func (tr TypeRecord) String() string {
	return [...]string{"Pairs login/password", "Text", "Binary", "Bank card"}[tr]
}

func (tr TypeRecord) Int() int {
	return [...]int{1, 2, 3, 4}[tr]
}

func (e EventDB) String() string {
	return [...]string{"edit", "del"}[e]
}
