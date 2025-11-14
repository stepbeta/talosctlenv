package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed talosctl versions",
	Long:  `Lists all the talosctl versions that are currently installed on the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vrsPath, err := utils.GetVrsPath(cmd)
		if err != nil {
			cmd.Println("Error getting vrs path:", err)
			return err
		}
		versions, err := utils.ListInstalledVersions(vrsPath)
		if err != nil {
			cmd.Println("Error listing available binaries:", err)
			return err
		}
		if len(versions) == 0 {
			cmd.Println("No talosctl versions installed.")
			return nil
		}

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

		cmd.Println("Available talosctl versions:")
		for _, v := range versions {
			vrs := v.Original()
			if vrs == currentVersion {
				vrs += " *"
			}
			cmd.Println(vrs)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
