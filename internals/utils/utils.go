package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// GetDefaultBinPath returns the default bin path for talosctl binaries.
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

// EnsurePathExists ensures that the given path exists, creating it if necessary.
func EnsurePathExists(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// osAliases returns a list of OS name aliases for the current Go OS name.
func GetOSAlias() string {
	goos := strings.ToLower(runtime.GOOS)
	switch goos {
	case "darwin":
		return "darwin"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

// archAliases returns a list of architecture name aliases for the current Go architecture name.
func GetArchAliases() []string {
	goarch := strings.ToLower(runtime.GOARCH)
	switch goarch {
	case "amd64":
		return []string{"amd64", "x86_64", "x86-64"}
	case "arm64":
		return []string{"arm64", "aarch64"}
	default:
		return []string{goarch}
	}
}
