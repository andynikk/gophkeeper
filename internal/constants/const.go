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
	TimeLivingCertificateYaer   = 10
	TimeLivingCertificateMounth = 0
	TimeLivingCertificateDay    = 0

	TypeEncryption = "sha512"

	AdressServer = "localhost:8080"

	QuerySelectUserWithWhereTemplate = `SELECT 
							* 
						FROM 
							gophkeeper."Users"
						WHERE 
							"User" = $1;`

	QuerySelectUserWithPasswordTemplate = `SELECT 
							* 
						FROM 
							gophkeeper."Users"
						WHERE 
							"User" = $1 and "Password" = $2;`

	QueryInsertUserTemplate = `INSERT INTO 
							gophkeeper."Users" ("User", "Password") 
						VALUES
							($1, $2);`

	QueryInsertPairsTemplate = `INSERT INTO gophkeeper."PairsLoginPassword"(
							"User", "TypePairs", "Name", "Password")
						VALUES ($1, $2, $3, $4);`

	QueryUpdatePairsTemplate = `UPDATE gophkeeper."PairsLoginPassword"
						SET "User"=$1, "TypePairs"=$2, "Name"=$3, "Password"=$4
						WHERE "User" = $1 and "TypePairs" = $2 and "Name" = $3;`

	QuerySelectPairsTemplate = `SELECT "User", "TypePairs", "Name", "Password" 
						FROM 
							gophkeeper."PairsLoginPassword"
						WHERE 
							"User" = $1;`

	QuerySelectOnePairsTemplate = `SELECT "User", "TypePairs", "Name", "Password" 
						FROM 
							gophkeeper."PairsLoginPassword"
						WHERE 
							"User" = $1 and "TypePairs" = $2 and "Name" = $3;`
	QueryDelOnePairsTemplate = `DELETE 
						FROM 
							gophkeeper."PairsLoginPassword"
						WHERE 
							"User" = $1 and "TypePairs" = $2 and "Name" = $3;`

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

	HeaderAuthorization  = "Authorization"
	HeaderMiddlewareBody = "Middleware-Body"
)

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
