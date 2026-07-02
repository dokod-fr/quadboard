package config

import "testing"

func TestDefaults(t *testing.T) {
	cfg := Defaults()

	if cfg.Server.Address != ":8080" {
		t.Fatalf("expected :8080, got %s", cfg.Server.Address)
	}

	if cfg.Logging.Level != "info" {
		t.Fatalf("expected info, got %s", cfg.Logging.Level)
	}

	if cfg.Theme.Name != "default" {
		t.Fatalf("expected default theme, got %s", cfg.Theme.Name)
	}
}

func TestLoad(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Address == "" {
		t.Fatal("server address should not be empty")
	}
}
