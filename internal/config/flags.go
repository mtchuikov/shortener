package config

import "github.com/spf13/pflag"

var (
	addr    = pflag.StringP("addr", "a", "127.0.0.1:8080/", "server address")
	baseURL = pflag.StringP("base", "b", "", "base url for shortened links")
	verbose = pflag.BoolP("verbose", "v", false, "print verbose logs")
)

func loadFromFlags(config *Config) {
	pflag.Parse()
	config.Addr = *addr
	config.BaseURL = *baseURL
	config.Verbose = *verbose
}
