package config

import (
	"github.com/spf13/pflag"
)

var (
	serverAddrFlag = pflag.StringP("addr", "a", "127.0.0.1:8080", "server address")
	baseURLFlag    = pflag.StringP("base", "b", "http://127.0.0.1:8080", "base url for shortened links")
	verboseFlag    = pflag.BoolP("verbose", "v", false, "print verbose logs")
)

func loadFromFlags(config *Config) {
	pflag.Parse()
	config.ServerAddr = *serverAddrFlag
	config.BaseURL = *baseURLFlag
	config.Verbose = *verboseFlag
}
