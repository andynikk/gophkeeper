// Package logger: логирование ошибок и информации
package logger

import (
	"github.com/rs/zerolog"
)

type Logger struct {
	Log zerolog.Logger
}

// ErrorLog логироание уровня ошибка
func (l *Logger) ErrorLog(err error) {
	l.Log.Error().Err(err).Msg("")
}

// InfoLog логироание уровня информация
func (l *Logger) InfoLog(infoString string) {
	l.Log.Info().Msgf("%+v\n", infoString)
}
