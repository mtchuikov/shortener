package config

import (
	"github.com/spf13/pflag"
)

var (
	serverAddrFlagDesc = "Specify the IP address and port for the server to listen on"
	serverAddrFlag     = pflag.StringP("server-addr", "a", "127.0.0.1:8080", serverAddrFlagDesc)

	baseURLFlagDesc = "Define the base URL used to generate shortened links"
	baseURLFlag     = pflag.StringP("base-url", "b", "http://127.0.0.1:8080/", baseURLFlagDesc)

	verboseFlagDesc = "Enable verbose logging at the debug level. Overrides the log-level flag to 'debug'"
	verboseFlag     = pflag.BoolP("verbose", "v", false, verboseFlagDesc)
)

func (c *Config) loadFromFlags() {
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	c.ServerAddr = *serverAddrFlag
	c.BaseURL = *baseURLFlag
	c.Verbose = *verboseFlag
}
