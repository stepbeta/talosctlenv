package github

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
)

var errReleaseNotFound = errors.New("release not found")

type GithubHelper struct {
	Client *github.Client
}

func New() GithubHelper {
	// Optional: Use token for higher rate limits:
	// - anonymous: 60 calls per hour
	// - authenticated: 5,000 calls per hour
    token := os.Getenv("GITHUB_TOKEN")
    client := github.NewClient(nil)
    if token != "" {
        client = client.WithAuthToken(token)
    }
	return GithubHelper{Client: client}
}

type FetchOptions struct {
    IncludeDevel bool
    Limit        int
}

func (gh *GithubHelper) FetchAllReleases(opts FetchOptions) ([]*github.RepositoryRelease, error) {
    ctx := context.Background()

    // TODO add progress indicator

    var allReleases []*github.RepositoryRelease
    page := 1

    // we use max possible value in order to limit occurrence of rate-limiting
    limit := 100
    if opts.Limit < 100 {
        limit = opts.Limit
    }

    for {
        releases, resp, err := gh.Client.Repositories.ListReleases(ctx, "siderolabs", "talos", &github.ListOptions{
            Page:    page,
            PerPage: limit,
        })
        if err != nil {
            return nil, err
        }

        for _, r := range releases {
            if opts.Limit > 0 && len(allReleases) >= opts.Limit {
                return allReleases, nil
            }
            allReleases = append(allReleases, r)
        }

        if resp.NextPage == 0 {
            break
        }
        page = resp.NextPage
    }

    return allReleases, nil
}

func (gh *GithubHelper) DownloadRelease(version, vrsPath string) error {
	ctx := context.Background()
	rel, _, err := gh.Client.Repositories.GetReleaseByTag(ctx, "siderolabs", "talos", version)
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
	rc, _, err := gh.Client.Repositories.DownloadReleaseAsset(ctx, "siderolabs", "talos", asset.GetID(), http.DefaultClient)
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
	destPath := filepath.Join(vrsPath, "talosctl-" + version)
	if err := os.Rename(tmpFile.Name(), destPath); err != nil {
		return fmt.Errorf("failed to move downloaded file to destination: %w", err)
	}

	// make it executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %w", err)
	}
	return nil
}