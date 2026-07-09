package config

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Providers ProvidersConfig `yaml:"providers"`
}

type ServerConfig struct {
	Address      string `yaml:"address"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ProvidersConfig struct {
	Quadlet QuadletConfig `yaml:"quadlet"`
}

type QuadletConfig struct {
	Paths []string `yaml:"paths"`
}
