package config

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Logging   LoggingConfig   `yaml:"logging"`
	Providers ProvidersConfig `yaml:"providers"`
	Auth      AuthConfig      `yaml:"auth"`
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

type AuthConfig struct {
	SecretKey string      `yaml:"secret_key"` // key to sign the session cookie HMAC
	Secure    bool        `yaml:"secure"`     // Be careful when setting this to false, as it will allow cookies to be sent over HTTP
	OIDC      *OIDCConfig `yaml:"oidc"`
}

type OIDCConfig struct {
	Issuer       string `yaml:"issuer"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}
