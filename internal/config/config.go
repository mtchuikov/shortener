package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

var conf config = config{serviceName: "shortener"}

type config struct {
	serviceName string
	serverAddr  string `env:"SERVER_ADDRESS"`
	baseURL     string `env:"BASE_URL"`
	databaseDSN string `env:"DATABASE_DSN"`
	fileStorage string `env:"FILE_STORAGE_PATH"`
	verbose     bool   `env:"VERBOSE"`
}

const wFailedToLoadConfigFromEnv = "failed to config from env: %w"

func Init() error {
	conf.loadFromFlags()

	err := env.Parse(&conf)
	if err != nil {
		return fmt.Errorf(wFailedToLoadConfigFromEnv, err)
	}

	hasSlash := strings.HasSuffix(conf.baseURL, "/")
	if !hasSlash {
		conf.baseURL = conf.baseURL + "/"
	}

	return nil
}

func ServiceName() string {
	return conf.serviceName
}

func ServerAddr() string {
	return conf.serverAddr
}

func BaseURL() string {
	return conf.baseURL
}

func DatabaseDSN() string {
	return conf.databaseDSN
}

func FileStorage() string {
	return conf.fileStorage
}

func Verbose() bool {
	return conf.verbose
}
