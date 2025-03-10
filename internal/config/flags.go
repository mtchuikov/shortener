package config

import (
	"strings"

	"github.com/spf13/pflag"
)

var (
	addrFlag    = pflag.StringP("addr", "a", "127.0.0.1:8080", "server address")
	baseURLFlag = pflag.StringP("base-url", "b", "http://127.0.0.1:8080/", "base url for shortened links")
)

func loadFromFlags(config *Config) {
	pflag.Parse()
	config.Addr = *addrFlag

	baseURL := *baseURLFlag
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	config.BaseURL = baseURL
}
