package cmd

import (
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "talosctlenv tool version",
	Long:  `Shows the current version of the talosctlenv tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("talosctlenv version v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
