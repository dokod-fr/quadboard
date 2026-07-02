package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the QuadBoard server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting QuadBoard...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
