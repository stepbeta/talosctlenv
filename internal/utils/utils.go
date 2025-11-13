package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// GetDefaultBinPath returns the default bin path for in-use talosctl binary.
func GetDefaultBinPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	binPath := filepath.Join(homeDir, ".talosctlenv", "bin")
	return binPath, nil
}

// GetBinPath retrieves the bin path from the command flags or returns the default path.
func GetBinPath(cmd *cobra.Command) (string, error) {
	binPath, err := cmd.Flags().GetString("bin-path")
	if err != nil {
		return "", err
	}
	if binPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		binPath = filepath.Join(homeDir, ".talosctlenv", "bin")
	}
	return binPath, nil
}

// GetDefaultVrsPath returns the default bin path for downloaded talosctl binaries.
func GetDefaultVrsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	binPath := filepath.Join(homeDir, ".talosctlenv", "versions")
	return binPath, nil
}

// GetBinPath retrieves the vrs path from the command flags or returns the default path.
func GetVrsPath(cmd *cobra.Command) (string, error) {
	vrsPath, err := cmd.Flags().GetString("vrs-path")
	if err != nil {
		return "", err
	}
	if vrsPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		vrsPath = filepath.Join(homeDir, ".talosctlenv", "versions")
	}
	return vrsPath, nil
}

// EnsurePathExists ensures that the given path exists, creating it if necessary.
func EnsurePathExists(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
