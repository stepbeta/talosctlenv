package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "generate talosctlenv documentation",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			cmd.Printf("Error generating docs: %v\n", err)
		} else {
			cmd.Println("Documentation generated in ./docs")
		}
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}