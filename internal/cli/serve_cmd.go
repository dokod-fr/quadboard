package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/auth"
	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/http"
	"github.com/dokod-fr/quadboard/internal/logging"
	"github.com/dokod-fr/quadboard/internal/providers/quadlet"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	RunE:  Run,
}

func Run(cmd *cobra.Command, args []string) error {
	cfg, configPath, err := config.Load()
	if err != nil {
		// Initialize a basic logger to report the config error
		logger := logging.New("info", "text")
		logger.Error("Failed to load configuration", "error", err)
		return err
	}

	logger := logging.New(cfg.Logging.Level, cfg.Logging.Format)
	slog.SetDefault(logger)

	logger.Info("Starting QuadBoard")

	if configPath != "" {
		logger.Info("Configuration file loaded successfully", "path", configPath)
	} else {
		logger.Info("No configuration file found, using defaults and environment variables")
	}

	logger.Info("Logging initialized", slog.String("level", cfg.Logging.Level))
	logger.Debug("Server configuration",
		slog.Int("read_timeout", cfg.Server.ReadTimeout),
		slog.Int("write_timeout", cfg.Server.WriteTimeout))
	logger.Info("Quadlet provider initialized", slog.Any("paths", cfg.Providers.Quadlet.Paths))

	// Initialize provider registry with Quadlet provider for now
	providerRegistry := app.NewProviderRegistry(
		quadlet.New(cfg.Providers.Quadlet.Paths...),
	)

	catalog := app.NewCatalog(providerRegistry)

	// First refresh at start
	logger.Info("Performing initial catalog scan...")
	if err := catalog.Refresh(); err != nil {
		logger.Warn("Initial catalog scan failed or completed with warnings", "error", err)
	} else {
		logger.Info("Catalog populated successfully", "resources_count", len(catalog.Resources()))
	}

	var oidcInstance *auth.OIDC
	if cfg.Auth.OIDC != nil && cfg.Auth.OIDC.Issuer != "" {
		logger.Info("Initializing OIDC authentication", "issuer", cfg.Auth.OIDC.Issuer)
		oidcInstance, err = auth.NewOIDC(
			context.Background(),
			cfg.Auth.OIDC.Issuer,
			cfg.Auth.OIDC.ClientID,
			cfg.Auth.OIDC.ClientSecret,
			cfg.BaseURL,
			cfg.Auth.SecretKey,
			cfg.Auth.Secure,
		)
		if err != nil {
			return fmt.Errorf("failed to initialize OIDC: %w", err)
		}
	} else {
		logger.Info("Authentication disabled. All resources will be visible.")
	}

	router := http.NewRouter(catalog, oidcInstance)

	server := http.NewServer(cfg, router)

	return server.Run()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
