package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internals/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed talosctl versions",
	Long:  `Lists all the talosctl versions that are currently installed on the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		binPath, err := utils.GetBinPath(cmd)
		if err != nil {
			cmd.Println("Error getting bin path:", err)
			return err
		}
		versions, err := listInstalledVersions(binPath)
		if err != nil {
			cmd.Println("Error listing available binaries:", err)
			return err
		}
		if len(versions) == 0 {
			cmd.Println("No talosctl versions installed.")
			return nil
		}
		cmd.Println("Available talosctl versions:")
		for _, v := range versions {
			cmd.Println(v)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listInstalledVersions(binPath string) ([]string, error) {
	files, err := os.ReadDir(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			// binPath does not exist, return empty list
			return []string{}, nil
		}
		return nil, err
	}
	versions := make([]string, 0)
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
		versions = append(versions, fv[1])
	}
	return versions, nil
}