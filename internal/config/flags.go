package config

import "github.com/spf13/pflag"

var (
	serverAddrFlagDesc = "Specify the IP address and port for the server to listen on"
	serverAddrFlag     = pflag.StringP("server-addr", "a", "127.0.0.1:8080", serverAddrFlagDesc)

	baseURLFlagDesc = "Define the base URL used to generate shortened links"
	baseURLFlag     = pflag.StringP("base-url", "b", "http://127.0.0.1:8080/", baseURLFlagDesc)

	databaseDSNDesc = "Define the database DSN (Data Source Name) used to connect to the database"
	databaseDSNFlag = pflag.StringP("dsn", "d", "postgres://user:password@127.0.0.1:5432/postgres?sslmode=disable", databaseDSNDesc)

	fileStorageDesc = "Path to the JSON file for storing shorten URLs and its identifiers"
	fileStorageFlag = pflag.StringP("file-storage", "f", "shorten-urls", fileStorageDesc)

	verboseFlagDesc = "Enable verbose logging at the debug level. Overrides the log-level flag to 'debug'"
	verboseFlag     = pflag.BoolP("verbose", "v", false, verboseFlagDesc)
)

func (c *config) loadFromFlags() {
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	c.serverAddr = *serverAddrFlag
	c.baseURL = *baseURLFlag
	c.databaseDSN = *databaseDSNFlag
	c.fileStorage = *fileStorageFlag
	c.verbose = *verboseFlag
}
