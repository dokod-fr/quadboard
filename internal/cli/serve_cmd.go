package cli

import (
	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/http"
	"github.com/dokod-fr/quadboard/internal/providers/quadlet"
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
			quadlet.New("internal/providers/quadlet/testdata"),
		)

		router := http.NewRouter(discovery)

		server := http.NewServer(cfg, router)

		return server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
