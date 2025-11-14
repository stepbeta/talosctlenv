package cmd

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internal/github"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

var includeDevel bool
var limit int

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List all remote talosctl versions from GitHub (sorted by semver)",
	Long: `List all talosctl versions available on GitHub, sorted by semver.

In the list the versions currently installed are marked with a '+' symbol, while the version currently in use is marked with a '*' symbol.

Note: By default pre-release versions (alpha, beta, rc) are hidden. Use '--devel' to include them`,
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

		// ignore errors here, it's not important
		currentVersion := ""
		binPath, err := utils.GetBinPath(cmd)
		if err != nil {
			cmd.Println("Error getting bin path:", err)
		}
		if binPath != "" {
			currentVersion, err = utils.GetVrsInUse(binPath)
			if err != nil {
				cmd.Println("Error getting current version:", err)
			}
		}

		// ignore errors here, it's not important
		var localVersions []*semver.Version
		vrsPath, err := utils.GetVrsPath(cmd)
		if err != nil {
			cmd.Println("Error getting vrs path:", err)
		}
		if vrsPath != "" {
			localVersions, err = utils.ListInstalledVersions(vrsPath)
			if err != nil {
				cmd.Println("Error listing available binaries:", err)
			}
		}
		cmd.Println("Available versions to download:")
		for _, v := range versions {
			vrs := v.Original()
			if vrs == currentVersion {
				vrs += " *"
			} else {
				for _, lv := range localVersions {
					if lv.Equal(v) {
						vrs += " +"
						break
					}
				}
			}
			cmd.Println(vrs)
		}

		if !includeDevel {
			cmd.Println("\nNote: Pre-release versions (alpha, beta, rc) are hidden. Use '--devel' to include them.")
		}
	},
}

func init() {
	listRemoteCmd.Flags().BoolVar(&includeDevel, "devel", false, "Include pre-release versions (alpha, beta, rc)")
	listRemoteCmd.Flags().IntVarP(&limit, "limit", "l", 0, "Limit number of versions displayed")
	rootCmd.AddCommand(listRemoteCmd)
}
