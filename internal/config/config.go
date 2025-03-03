package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr    string `env:"SERVER_ADDRESS"`
	BaseURL string `env:"BASE_URL"`
	Verbose bool   `env:"VERBOSE"`
}

func New() (Config, error) {
	config := Config{}
	loadFromFlags(&config)

	err := env.Parse(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
