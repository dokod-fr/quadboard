package cli

import (
	"context"

	"github.com/dokod-fr/quadboard/internal/config"
	httpserver "github.com/dokod-fr/quadboard/internal/http"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		srv := httpserver.New(cfg)

		return srv.Start(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
