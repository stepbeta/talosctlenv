package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

var (
	errFileNotFound = errors.New("version not found")

	useCmd = &cobra.Command{
		Use:   "use <version>",
		Short: "Set the specified talosctl version as the active one",
		Long: `Create a symlink to the specified version with the name "talosctl".

Make sure the "bin-path" is included in the $PATH variable.`,
		Args:  cobra.ExactArgs(1),
		RunE: useVersion,
	}
)

func useVersion(cmd *cobra.Command, args []string) error {
	binPath, err := utils.GetBinPath(cmd)
	if err != nil {
		cmd.Println("Error getting bin path:", err)
		return err
	}
	err = utils.EnsurePathExists(binPath)
	if err != nil {
		cmd.Println("Error ensuring bin path exists:", err)
		return err
	}
	vrsPath, err := utils.GetVrsPath(cmd)
	if err != nil {
		cmd.Println("Error getting vrs path:", err)
		return err
	}
	fileName := filepath.Join(vrsPath, "talosctl-" + args[0])
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		install, err := cmd.Flags().GetBool("install")
		if err != nil {
			cmd.Println("Error retrieving 'install' flag:", err)
			cmd.Println("Skipping action")
			return err
		}
		if !install {
			cmd.Println("Error: specified version is not installed. Please install it first using `talosctlenv install <version>`")
			return errFileNotFound
		}
		// here we want to install the version if not present
		rootCmd.SetArgs([]string{"install", args[0]})
		if err := rootCmd.Execute(); err != nil {
			cmd.Println("Error executing install:", err)
			cmd.Println("Skipping action")
			return err
		}
		// here we should have installed the version, we assume it succeeded
	}
	target := filepath.Join(binPath, "talosctl")
	// Check if the symlink already exists
	if _, err := os.Lstat(target); err == nil {
		// Remove existing symlink or file
		if err := os.Remove(target); err != nil {
			cmd.Println("Error removing existing symlink:", err)
			return err
		}
	}
	// create new symlink
	if err := os.Symlink(fileName, target); err != nil {
		cmd.Println("Error creating symlink:", err)
		return err
	}

	cmd.Printf("Now using talosctl version %s\n", args[0])
	return nil
}

func init() {
	useCmd.Flags().BoolP("install", "i", false, "Install the version if not yet present (best effort)")
	rootCmd.AddCommand(useCmd)
}
