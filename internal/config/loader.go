package config

import (
	"fmt"
	"log/slog"
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

	slog.Debug("Configuration defined", slog.Any("config", cfg))
	return cfg, configPath, nil
}

// resolveConfigPath determines which configuration file to use.
// 1. Uses QUADBOARD_CONFIG_FILE if set.
// 2. Looks for config.yaml next to the executable.
// 3. Returns empty string if no file is found.
func resolveConfigPath() string {
	if path := os.Getenv("QUADBOARD_CONFIG_FILE"); path != "" {
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
	if val, ok := os.LookupEnv("QUADBOARD_BASE_URL"); ok {
		cfg.BaseURL = strings.TrimSuffix(val, "/") // Optionnel : nettoie le slash final
	}
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

	// --- UI --- //
	if val, ok := os.LookupEnv("QUADBOARD_UI_SHOW_ITSELF"); ok {
		if show, err := strconv.ParseBool(val); err == nil {
			cfg.UI.ShowItSelf = show
		}
	}

	if val, ok := os.LookupEnv("QUADBOARD_UI_HIDE_UNAUTHORIZED_GROUPS"); ok {
		if hide, err := strconv.ParseBool(val); err == nil {
			cfg.UI.HideUnauthorizedGroups = hide
		}
	}

	if val, ok := os.LookupEnv("QUADBOARD_UI_GROUP_LABELS"); ok {
		if cfg.UI.GroupLabels == nil {
			cfg.UI.GroupLabels = make(map[string]string)
		}
		// Expected format: "security:Sécurité,infra:Infrastructure"
		pairs := strings.Split(val, ",")
		for _, p := range pairs {
			kv := strings.SplitN(p, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				cfg.UI.GroupLabels[key] = value
			}
		}
	}

	if val, ok := os.LookupEnv("QUADBOARD_UI_GROUP_BY_DEFAULT"); ok {
		if b, err := strconv.ParseBool(val); err == nil {
			cfg.UI.GroupByDefault = b
		}
	}

	// ----- AUTH ----- //
	if val, ok := os.LookupEnv("QUADBOARD_AUTH_SECRET_KEY"); ok {
		cfg.Auth.SecretKey = val
	}
	if val, ok := os.LookupEnv("QUADBOARD_AUTH_SECURE"); ok {
		if secure, err := strconv.ParseBool(val); err == nil {
			cfg.Auth.Secure = secure
		}
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
	}
}
