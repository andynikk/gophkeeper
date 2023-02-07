package errs

import (
	"errors"
	"net/http"
)

// InvalidFormat Ошибка в тексте SQL запроса
var InvalidFormat = errors.New("invalid format")

// ErrLoginBusy ошибка создания пользователя. Имя занято.
var ErrLoginBusy = errors.New("login busy")

// ErrErrorServer ошибка на сервере.
var ErrErrorServer = errors.New("error server")

// ErrInvalidLoginPassword пара пользователь и пароль не найдены.
var ErrInvalidLoginPassword = errors.New("invalid login password")

// HTTPErrors Приведение ошибки к HTTP статусам
func HTTPErrors(err error) int {

	HTTPAnswer := http.StatusOK

	if errors.Is(err, InvalidFormat) {
		HTTPAnswer = http.StatusBadRequest
	} else if errors.Is(err, ErrLoginBusy) {
		HTTPAnswer = http.StatusConflict
	} else if errors.Is(err, ErrErrorServer) {
		HTTPAnswer = http.StatusInternalServerError
	} else if errors.Is(err, ErrInvalidLoginPassword) {
		HTTPAnswer = http.StatusUnauthorized
	}
	return HTTPAnswer
}
