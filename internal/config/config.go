package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

var conf config = config{ServiceName: "shortener"}

type config struct {
	ServiceName string
	ServerAddr  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	Verbose     bool   `env:"VERBOSE"`
}

const wFailedToLoadConfigFromEnv = "failed to config from env: %w"

func Init() error {
	conf.loadFromFlags()

	err := env.Parse(&conf)
	if err != nil {
		return fmt.Errorf(wFailedToLoadConfigFromEnv, err)
	}

	hasSlash := strings.HasSuffix(conf.BaseURL, "/")
	if !hasSlash {
		conf.BaseURL = conf.BaseURL + "/"
	}

	return nil
}

func ServiceName() string {
	return conf.ServiceName
}

func ServerAddr() string {
	return conf.ServerAddr
}

func BaseURL() string {
	return conf.BaseURL
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
