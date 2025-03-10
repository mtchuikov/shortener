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

const serviceName = "shortener"

func New() Config {
	config := Config{ServiceName: serviceName}
	loadFromFlags(&config)

	env.Parse(&config)

	if !strings.HasSuffix(config.BaseURL, "/") {
		config.BaseURL = config.BaseURL + "/"
	}

	return config
}
