package cmd

import (
	"fmt"

	"github.com/dokod-fr/quadboard/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`QuadBoard

Version : %s
Commit  : %s
Built   : %s
`,
			version.Version,
			version.Commit,
			version.Date,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
