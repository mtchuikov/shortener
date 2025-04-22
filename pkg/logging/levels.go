package logging

import (
	"errors"

	"github.com/rs/zerolog"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
	Panic
)

var ErrInvalidLogLevel = errors.New("invalid log level")

func NewLevel(level string) (Level, error) {
	switch level {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	case "fatal":
		return Fatal, nil
	case "panic":
		return Panic, nil
	default:
		return -1, ErrInvalidLogLevel
	}
}

func (l Level) ToZerolog() (zerolog.Level, error) {
	switch l {
	case Debug:
		return zerolog.DebugLevel, nil
	case Info:
		return zerolog.InfoLevel, nil
	case Warn:
		return zerolog.WarnLevel, nil
	case Error:
		return zerolog.ErrorLevel, nil
	case Fatal:
		return zerolog.FatalLevel, nil
	case Panic:
		return zerolog.PanicLevel, nil
	default:
		return -1, ErrInvalidLogLevel
	}
}
