// Package token: работа с токеном jwt
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"gophkeeper/internal/constants"
)

type Claims struct {
	Authorized bool
	User       string
	Exp        int64
}

// GenerateJWT генерация токена для пользователя
func (c *Claims) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = c.Authorized
	claims["user"] = c.User
	claims["exp"] = c.Exp

	tokenString, err := token.SignedString(constants.HashKey)

	if err != nil {
		constants.Logger.ErrorLog(err)
		return "", err
	}

	return tokenString, nil
}

// ExtractClaims получение имя пользователя из токена
func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := constants.HashKey
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		constants.Logger.InfoLog("Invalid JWT Token")
		return nil, false
	}
}

// NewClaims Инициализация сущности Claims
func NewClaims(name string) *Claims {
	return &Claims{
		Authorized: true,
		User:       name,
		Exp:        time.Now().Add(time.Hour * constants.TimeLiveToken).Unix(),
	}
}
