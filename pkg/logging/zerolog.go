package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZerolog(appName string) zerolog.Logger {
	zerolog.LevelFieldName = LevelFieldName
	zerolog.ErrorFieldName = ErrorFieldName
	zerolog.MessageFieldName = MessageFieldName
	zerolog.TimeFieldFormat = TimeFieldFormat

	log := zerolog.New(os.Stdout).With().
		Str("app", appName).Timestamp().
		Logger()

	return log.Level(zerolog.InfoLevel)
}
