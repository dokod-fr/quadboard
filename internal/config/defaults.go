package config

func Defaults() Config {
	return Config{
		Server: ServerConfig{
			Address: ":8080",
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
