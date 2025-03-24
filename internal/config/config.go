package config

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServiceName string
	ServerAddr  string `env:"SERVER_ADDRESS"`
	BaseURL     string `env:"BASE_URL"`
	Verbose     bool   `env:"VERBOSE"`
}

func New() Config {
	const serviceName = "shortener"
	config := Config{ServiceName: serviceName}

	config.loadFromFlags()
	env.Parse(&config)

	hasSlash := strings.HasSuffix(config.BaseURL, "/")
	if !hasSlash {
		config.BaseURL = config.BaseURL + "/"
	}

	return config
}
