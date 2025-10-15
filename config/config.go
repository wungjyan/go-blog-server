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
	Log struct {
		Level      string
		FilePath   string
		MaxSize    int
		MaxBackups int
		MaxAge     int
		Compress   bool
		Format     string
		ShowCaller bool
		Console    bool
	}
}
