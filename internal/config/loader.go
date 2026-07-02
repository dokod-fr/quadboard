package config

import (
	"fmt"
)

func Load() (Config, error) {
	cfg := Defaults()

	// TODO PR #3:
	// - TOML file
	// - env override
	// - flags override

	if err := validate(cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}
