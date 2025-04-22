package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

var conf config

type config struct {
	LogToFile    bool   `ENV:"LOG_TO_FILE"`
	LogFile      string `ENV:"LOG_FILE"`
	LogFileLevel string `ENV:"LOG_FILE_LEVEL"`
	ServerAddr   string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	JWTSecret    string `env:"JWT_SECRET"`
	DatabaseDSN  string `ENV:"DATABASE_DSN"`
	FileStorage  string `env:"FILE_STORAGE_PATH"`
	Verbose      bool   `env:"VERBOSE"`
}

var ErrUnableToParseEnv = errors.New("unable to parse env")

func Init() error {
	initFromFlags()

	err := env.Parse(&conf)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnableToParseEnv, err)
	}

	if !strings.HasSuffix(conf.BaseURL, "/") {
		conf.BaseURL = conf.BaseURL + "/"
	}

	return nil
}

func LogToFile() bool {
	return conf.LogToFile
}

func LogFile() string {
	return conf.LogFile
}

func LogFileLevel() string {
	return conf.LogFileLevel
}

func ServerAddr() string {
	return conf.ServerAddr
}

func BaseURL() string {
	return conf.BaseURL
}

func JWTSecret() string {
	return conf.JWTSecret
}

func DatabaseDSN() string {
	return conf.DatabaseDSN
}

func FileStorage() string {
	return conf.FileStorage
}

func Verbose() bool {
	return conf.Verbose
}
