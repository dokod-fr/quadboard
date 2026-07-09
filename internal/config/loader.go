package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load loads the configuration respecting the following priority:
// Environment variables > YAML file > Default constants.
func Load() (Config, error) {
	cfg := defaultConfig()

	if configPath := os.Getenv("QUADBOARD_CONFIG_PATH"); configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return cfg, fmt.Errorf("failed to read configuration file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, fmt.Errorf("failed to parse YAML file: %w", err)
		}
	}

	applyEnvVars(&cfg)

	return cfg, nil
}

func applyEnvVars(cfg *Config) {
	if val, ok := os.LookupEnv("QUADBOARD_SERVER_ADDRESS"); ok {
		cfg.Server.Address = val
	}
	if val, ok := os.LookupEnv("QUADBOARD_SERVER_READ_TIMEOUT"); ok {
		if t, err := strconv.Atoi(val); err == nil {
			cfg.Server.ReadTimeout = t
		}
	}
	if val, ok := os.LookupEnv("QUADBOARD_SERVER_WRITE_TIMEOUT"); ok {
		if t, err := strconv.Atoi(val); err == nil {
			cfg.Server.WriteTimeout = t
		}
	}

	if val, ok := os.LookupEnv("QUADBOARD_LOGGING_LEVEL"); ok {
		cfg.Logging.Level = val
	}
	if val, ok := os.LookupEnv("QUADBOARD_LOGGING_FORMAT"); ok {
		cfg.Logging.Format = val
	}

	if val, ok := os.LookupEnv("QUADBOARD_QUADLET_PATHS"); ok {
		paths := strings.Split(val, ",")
		for i := range paths {
			paths[i] = strings.TrimSpace(paths[i])
		}
		cfg.Providers.Quadlet.Paths = paths
	}
}
