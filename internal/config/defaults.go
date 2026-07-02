package config

func Defaults() Config {
	return Config{
		Server: ServerConfig{
			Address:      ":8080",
			ReadTimeout:  5,
			WriteTimeout: 10,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Theme: ThemeConfig{
			Name: "default",
		},
	}
}
