package config

type Config struct {
	ServiceName string
	Addr        string
	BaseURL     string
}

const serviceName = "shortener"

func New() Config {
	config := Config{ServiceName: serviceName}
	loadFromFlags(&config)
	return config
}
