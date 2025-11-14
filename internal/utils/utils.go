package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
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

// ListInstalledVersions lists all installed talosctl versions in the given vrsPath.
func ListInstalledVersions(vrsPath string) ([]*semver.Version, error) {
	files, err := os.ReadDir(vrsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// vrsPath does not exist, return empty list
			return []*semver.Version{}, nil
		}
		return nil, err
	}
	versions := make([]*semver.Version, 0)
	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() || !strings.HasPrefix(fileName, "talosctl") {
			// skip directories and non-talosctl files
			continue
		}
		// by convention the file name is talosctl-VERSION
		fv := strings.Split(fileName, "-")
		if fv == nil || len(fv) != 2 {
			// skip unexpected file names
			continue
		}
		v, err := semver.NewVersion(fv[1])
		if err == nil {
			versions = append(versions, v)
		}
	}
	sort.Sort(semver.Collection(versions))

	return versions, nil
}

// GetVrsInUse returns the version of talosctl currently in use from the given binPath.
func GetVrsInUse(binPath string) (string, error) {
	linkPath, err := filepath.EvalSymlinks(filepath.Join(binPath, "talosctl"))
	if err != nil {
		return "", err
	}
	if linkPath == "" {
		return "", nil
	}
	baseName := filepath.Base(linkPath)
	parts := strings.Split(baseName, "-")
	if len(parts) == 2 {
		return parts[1], nil
	}
	return "", nil
}
