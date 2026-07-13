package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	envKeys := []string{
		"QUADBOARD_CONFIG_PATH",
		"QUADBOARD_SERVER_ADDRESS",
		"QUADBOARD_SERVER_READ_TIMEOUT",
		"QUADBOARD_SERVER_WRITE_TIMEOUT",
		"QUADBOARD_LOGGING_LEVEL",
		"QUADBOARD_LOGGING_FORMAT",
		"QUADBOARD_QUADLET_PATHS",
	}
	for _, key := range envKeys {
		defer os.Unsetenv(key)
	}

	t.Run("Default only", func(t *testing.T) {
		os.Unsetenv("QUADBOARD_CONFIG_PATH")
		cfg, path, err := Load()
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		if path != "" {
			t.Errorf("Expected empty path, got: %s", path)
		}

		if cfg.Server.Address != "0.0.0.0:8080" {
			t.Errorf("Expected default address: 0.0.0.0:8080, got: %s", cfg.Server.Address)
		}
		if cfg.Logging.Level != "info" {
			t.Errorf("Expected default level: info, got: %s", cfg.Logging.Level)
		}
		if len(cfg.Providers.Quadlet.Paths) != 2 {
			t.Errorf("Expected 2 default paths, got: %v", cfg.Providers.Quadlet.Paths)
		}
	})

	t.Run("YAML overrides defaults", func(t *testing.T) {
		dir := t.TempDir()
		yamlPath := filepath.Join(dir, "config.yaml")
		yamlContent := []byte(`
server:
  address: "127.0.0.1:9090"
  read_timeout: 15
logging:
  level: "debug"
providers:
  quadlet:
    paths:
      - /custom/quadlet/path
`)
		if err := os.WriteFile(yamlPath, yamlContent, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		os.Setenv("QUADBOARD_CONFIG_PATH", yamlPath)
		os.Unsetenv("QUADBOARD_SERVER_ADDRESS")
		os.Unsetenv("QUADBOARD_LOGGING_LEVEL")

		cfg, path, err := Load()
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}
		if path != yamlPath {
			t.Errorf("Expected path %s, got: %s", yamlPath, path)
		}

		if cfg.Server.Address != "127.0.0.1:9090" {
			t.Errorf("Expected YAML address: 127.0.0.1:9090, got: %s", cfg.Server.Address)
		}
		if cfg.Logging.Level != "debug" {
			t.Errorf("Expected YAML level: debug, got: %s", cfg.Logging.Level)
		}
		if len(cfg.Providers.Quadlet.Paths) != 1 || cfg.Providers.Quadlet.Paths[0] != "/custom/quadlet/path" {
			t.Errorf("Invalid YAML paths: %v", cfg.Providers.Quadlet.Paths)
		}
	})

	t.Run("EnvVar overrides YAML", func(t *testing.T) {
		dir := t.TempDir()
		yamlPath := filepath.Join(dir, "config.yaml")
		yamlContent := []byte(`
server:
  address: "127.0.0.1:9090"
providers:
  quadlet:
    paths:
      - /etc/quadlet
`)
		if err := os.WriteFile(yamlPath, yamlContent, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		os.Setenv("QUADBOARD_CONFIG_PATH", yamlPath)
		os.Setenv("QUADBOARD_SERVER_ADDRESS", "0.0.0.0:3000")
		os.Setenv("QUADBOARD_LOGGING_LEVEL", "error")
		os.Setenv("QUADBOARD_QUADLET_PATHS", "/var/run/quadlet, /opt/quadlet")

		cfg, _, err := Load()
		if err != nil {
			t.Fatalf("Load() unexpected error: %v", err)
		}

		if cfg.Server.Address != "0.0.0.0:3000" {
			t.Errorf("Expected EnvVar address: 0.0.0.0:3000, got: %s", cfg.Server.Address)
		}
		if cfg.Logging.Level != "error" {
			t.Errorf("Expected EnvVar level: error, got: %s", cfg.Logging.Level)
		}
		if len(cfg.Providers.Quadlet.Paths) != 2 || cfg.Providers.Quadlet.Paths[1] != "/opt/quadlet" {
			t.Errorf("Invalid EnvVar paths: %v", cfg.Providers.Quadlet.Paths)
		}
	})
}
