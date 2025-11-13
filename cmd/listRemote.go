package cmd

import (
    "fmt"
    "sort"

    "github.com/Masterminds/semver/v3"
    "github.com/spf13/cobra"
    "github.com/stepbeta/talosctlenv/internal/github"
)

var includeDevel bool
var limit int

var listRemoteCmd = &cobra.Command{
    Use:   "list-remote",
    Short: "List all remote talosctl versions from GitHub (sorted by semver)",
    Run: func(cmd *cobra.Command, args []string) {
		ghc := github.New()
        releases, err := ghc.FetchAllReleases(github.FetchOptions{
            IncludeDevel: includeDevel,
            Limit:        limit,
        })
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        var versions []*semver.Version
        for _, r := range releases {
            v, err := semver.NewVersion(*r.TagName)
            if err == nil {
                if !includeDevel && v.Prerelease() != "" {
                    continue
                }
                versions = append(versions, v)
            }
        }

        sort.Sort(semver.Collection(versions))

        for _, v := range versions {
			// TODO if v is currently in use add a symbol
			// TODO if v is not currently in use but it's available locally, add a different symbol
            fmt.Println(v.Original())
        }


        if !includeDevel {
            fmt.Println("\nNote: Pre-release versions (alpha, beta, rc) are hidden. Use '--devel' to include them.")
        }
    },
}

func init() {
    listRemoteCmd.Flags().BoolVar(&includeDevel, "devel", false, "Include pre-release versions (alpha, beta, rc)")
    listRemoteCmd.Flags().IntVar(&limit, "limit", 0, "Limit number of versions displayed")
    rootCmd.AddCommand(listRemoteCmd)
}
