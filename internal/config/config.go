package config

type Config struct {
	Addr    string
	BaseURL string
	Verbose bool
}

func New() Config {
	config := Config{}
	loadFromFlags(&config)
	return config
}
