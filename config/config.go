package config

type Config struct {
	App struct {
		Port int
	}
	DB struct {
		Host     string
		Port     int
		Database string
		Username string
		Password string
	}
}
