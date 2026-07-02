package config

import "errors"

func validate(cfg Config) error {
	if cfg.Server.Address == "" {
		return errors.New("server address cannot be empty")
	}
	return nil
}
