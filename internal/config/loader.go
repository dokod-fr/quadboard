package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load loads the configuration respecting the following priority:
// Environment variables > YAML file > Default constants.
// It returns the Config, the path of the config file used (if any), and an error.
func Load() (Config, string, error) {
	cfg := defaultConfig()

	configPath := resolveConfigPath()
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return cfg, configPath, fmt.Errorf("failed to read configuration file: %w", err)
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return cfg, configPath, fmt.Errorf("failed to parse YAML file: %w", err)
		}
	}

	applyEnvVars(&cfg)

	return cfg, configPath, nil
}

// resolveConfigPath determines which configuration file to use.
// 1. Uses QUADBOARD_CONFIG_PATH if set.
// 2. Looks for config.yaml next to the executable.
// 3. Returns empty string if no file is found.
func resolveConfigPath() string {
	if path := os.Getenv("QUADBOARD_CONFIG_PATH"); path != "" {
		return path
	}

	exePath, err := os.Executable()
	if err != nil {
		return ""
	}

	fmt.Println("Executable path:", exePath) // Debugging line

	adjacentPath := filepath.Join(filepath.Dir(exePath), "config.yaml")
	if _, err := os.Stat(adjacentPath); err == nil {
		return adjacentPath
	}

	return ""
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

	if val, ok := os.LookupEnv("QUADBOARD_AUTH_SECRET_KEY"); ok {
		cfg.Auth.SecretKey = val
	}

	// Initialize OIDCConfig pointer if any OIDC-related environment variable is set
	if cfg.Auth.OIDC == nil {
		if _, ok := os.LookupEnv("QUADBOARD_AUTH_OIDC_ISSUER"); ok {
			cfg.Auth.OIDC = &OIDCConfig{}
		}
	}

	if cfg.Auth.OIDC != nil {
		if val, ok := os.LookupEnv("QUADBOARD_AUTH_OIDC_ISSUER"); ok {
			cfg.Auth.OIDC.Issuer = val
		}
		if val, ok := os.LookupEnv("QUADBOARD_AUTH_OIDC_CLIENT_ID"); ok {
			cfg.Auth.OIDC.ClientID = val
		}
		if val, ok := os.LookupEnv("QUADBOARD_AUTH_OIDC_CLIENT_SECRET"); ok {
			cfg.Auth.OIDC.ClientSecret = val
		}
		if val, ok := os.LookupEnv("QUADBOARD_AUTH_OIDC_REDIRECT_URL"); ok {
			cfg.Auth.OIDC.RedirectURL = val
		}
	}
}
