package config

import "github.com/spf13/pflag"

func initFromFlags() {
	pflag.BoolVarP(
		&conf.LogToFile,
		"log.file.enable", "",
		false,
		"Flag to switch logging to file",
	)

	pflag.StringVarP(
		&conf.LogFile,
		"log.file", "",
		"shortener.log",
		"File to write logs",
	)

	pflag.StringVarP(
		&conf.LogFileLevel,
		"log.file.level", "",
		"error",
		"Min log level to write to file",
	)

	pflag.StringVarP(
		&conf.ServerAddr,
		"addr", "a",
		"127.0.0.1:8080",
		"Address to listen http server",
	)

	pflag.StringVarP(
		&conf.BaseURL,
		"base", "b",
		"http://127.0.0.1:8080/",
		"Base for shorten URLs. Shortened URL = base + random id",
	)

	pflag.StringVarP(
		&conf.JWTSecret,
		"secret", "s",
		"jwtsecret", "Secret to sign and verify JWTs",
	)

	pflag.StringVarP(
		&conf.DatabaseDSN,
		"dsn", "d",
		"", "URL to connect Postgres database",
	)

	pflag.StringVarP(
		&conf.FileStorage,
		"file", "f",
		"shortener.backup",
		"File to backup shorten URLs from cache",
	)

	pflag.BoolVarP(
		&conf.Verbose,
		"verbose", "v",
		false,
		"Flag to switch verbose mode of logging",
	)

	pflag.CommandLine.SortFlags = false
	pflag.Parse()
}
