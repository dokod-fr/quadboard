package cli

import (
	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/http"
	"github.com/dokod-fr/quadboard/internal/mock"
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

		discovery := app.NewDiscovery(
			mock.New(),
		)

		router := http.NewRouter(discovery)

		server := http.NewServer(cfg, router)

		return server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
