package config

func defaultConfig() Config {
	quadletPaths := []string{
		"/etc/containers/systemd/",
	}

	return Config{
		Server: ServerConfig{
			Address:      "0.0.0.0:8080",
			ReadTimeout:  5,
			WriteTimeout: 10,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Providers: ProvidersConfig{
			Quadlet: QuadletConfig{
				Paths: quadletPaths,
			},
		},
		Auth: AuthConfig{
			SecretKey: "",
			Secure:    true, // Default to true for security; can be overridden for testing purposes
			OIDC:      nil,
		},
		UI: UIConfig{
			ShowItSelf:             false, // Default, no show quadboard
			GroupLabels:            map[string]string{},
			HideUnauthorizedGroups: false, // Default, show unauthorized group desactived
			GroupByDefault:         false, // Default, show flate view
		},
	}
}
