package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internals/utils"
)

var (
	errFileNotFound = errors.New("version not found")

	useCmd = &cobra.Command{
		Use:   "use <version>",
		Short: "Set the specified talosctl version as the active one",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			binPath, err := utils.GetBinPath(cmd)
			if err != nil {
				cmd.Println("Error getting bin path:", err)
				return err
			}
			fileName := filepath.Join(binPath, "talosctl-" + args[0])
			if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
				cmd.Println("Error: specified version is not installed. Please install it first using `talosctlenv install <version>`")
				return errFileNotFound
			}

			if err := os.Symlink(fileName, filepath.Join(binPath, "talosctl")); err != nil {
				cmd.Println("Error creating symlink:", err)
				return err
			}

			cmd.Printf("Now using talosctl version %s\n", args[0])
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(useCmd)
}
