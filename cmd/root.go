package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

// rootCmd represents the base command when called without any subcommands
var (
	// override at build time using `go build -ldflags "-X github.com/stepbeta/talosctlenv/cmd.Version=x.y.z"`
	Version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:   "talosctlenv",
		Short: "A talosctl versions manager",
		Long: `A tool to easily install and use multiple versions of talosctl.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	defaultBinPath, err := utils.GetDefaultBinPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error determining default bin path:", err)
		os.Exit(1)
	}
	rootCmd.PersistentFlags().StringP("bin-path", "b", defaultBinPath, "Absolute path to folder storing talosctl binaries")
}
