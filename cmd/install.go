package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v78/github"
	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internals/utils"
)

var (
	errReleaseNotFound = errors.New("release not found")

	installCmd = &cobra.Command{
		Use:   "install <version>",
		Short: "Download and install talosctl for the current OS/ARCH",
		Args:  cobra.ExactArgs(1),
		RunE: installVersion,
	}
)

func init() {
	rootCmd.AddCommand(installCmd)
}

func installVersion(cmd *cobra.Command, args []string) error {
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
	// TODO check if the version is already installed

	gh := github.NewClient(nil)
	if os.Getenv("GITHUB_TOKEN") != "" {
		gh = gh.WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	}
	ctx := context.Background()
	rel, _, err := gh.Repositories.GetReleaseByTag(ctx, "siderolabs", "talos", args[0])
	if err != nil {
		return err
	}
	if rel == nil {
		return errReleaseNotFound
	}
	osAlias := strings.ToLower(runtime.GOOS)
	archAlias := strings.ToLower(runtime.GOARCH)

	relName := "talosctl-" + osAlias + "-" + archAlias
	var asset *github.ReleaseAsset
	for _, a := range rel.Assets {
		if a == nil {
			continue
		}
		lname := strings.ToLower(a.GetName())
		if !strings.HasPrefix(lname, relName) {
			// not the right asset
			continue
		}
		if osAlias == "windows" && !strings.HasSuffix(lname, ".exe") {
			// windows binary must have .exe suffix
			continue
		}
		if osAlias != "windows" && len(strings.Split(lname, ".")) > 1 {
			// non-windows binaries should not have an extension
			continue
		}
		asset = a
	}
	if asset == nil {
		return errReleaseNotFound
	}

	// download asset using go-github helper (returns ReadCloser)
	rc, _, err := gh.Repositories.DownloadReleaseAsset(ctx, "siderolabs", "talos", asset.GetID(), http.DefaultClient)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer rc.Close()
	
	// write to temp file then move (safer)
	tmpFile, err := os.CreateTemp("", "talosctl-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	_, err = io.Copy(tmpFile, rc)
	if err1 := tmpFile.Close(); err == nil && err1 != nil {
		err = err1
	}
	if err != nil {
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to save download: %w", err)
	}

	// move to destination
	destPath := filepath.Join(binPath, "talosctl-" + args[0])
	if err := os.Rename(tmpFile.Name(), destPath); err != nil {
		return fmt.Errorf("failed to move downloaded file to destination: %w", err)
	}

	// make executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}

	cmd.Printf("talosctl version %s successfully installed\n", args[0])
	cmd.Printf("To switch to that version run `talosctlenv use %s`\n", args[0])

	return nil
}
