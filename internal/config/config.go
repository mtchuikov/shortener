package config

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServiceName string
	ServerAddr  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	FileStorage string `env:"FILE_STORAGE_PATH"`
	Verbose     bool   `env:"VERBOSE"`
}

func New() Config {
	config := Config{ServiceName: "shortener"}

	config.loadFromFlags()
	env.Parse(&config)

	hasSlash := strings.HasSuffix(config.BaseURL, "/")
	if !hasSlash {
		config.BaseURL = config.BaseURL + "/"
	}

	return config
}
