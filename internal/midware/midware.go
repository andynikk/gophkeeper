// Package midware: middleware сервера. Проверка на аутетификацию
package midware

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"gophkeeper/internal/constants"
)

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")

		if r.Header["Authorization"] != nil {

			TokenFindMatches(endpoint, w, r)
			return
		}
		TokenNotFound(w)
	})
}

func TokenFindMatches(endpoint func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("there was an error")
		}
		return constants.HashKey, nil
	})
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		return
	}

	if token.Valid {
		endpoint(w, r)
	}
}

func TokenNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("Not Authorized"))
	if err != nil {
		constants.Logger.ErrorLog(err)
	}
}
