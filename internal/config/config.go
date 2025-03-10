package config

type Config struct {
	ServiceName string
	Addr        string
	BaseURL     string
	Verbose     bool
}

const serviceName = "shortener"

func New() Config {
	config := Config{ServiceName: serviceName}
	loadFromFlags(&config)
	return config
}
