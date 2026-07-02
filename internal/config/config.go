package config

type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
	Theme   ThemeConfig
}

type ServerConfig struct {
	Address string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type ThemeConfig struct {
	Name string
}
