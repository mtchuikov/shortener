package config

import (
	"strings"

	"github.com/spf13/pflag"
)

var (
	addrFlag    = pflag.StringP("addr", "a", "127.0.0.1:8080", "server address")
	baseURLFlag = pflag.StringP("base-url", "b", "http://127.0.0.1:8080/", "base url for shortened links")
	verboseFlag = pflag.BoolP("verbose", "v", false, "print verbose logs")
)

func loadFromFlags(config *Config) {
	pflag.Parse()
	config.Addr = *addrFlag

	baseURL := *baseURLFlag
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	config.BaseURL = baseURL
	config.Verbose = *verboseFlag
}
