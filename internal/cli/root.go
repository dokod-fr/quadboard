package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "quadboard",
	Short: "QuadBoard is a zero-config application portal for Podman Quadlets",
	Long: `QuadBoard automatically discovers applications deployed with
Podman Quadlets and exposes them through a clean web interface.`,
}

func Execute() error {
	return rootCmd.Execute()
}
