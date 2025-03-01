package config

type Config struct {
	Addr string
}

func New() Config {
	return Config{
		Addr: "127.0.0.1:8080",
	}
}
