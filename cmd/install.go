package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stepbeta/talosctlenv/internal/github"
	"github.com/stepbeta/talosctlenv/internal/utils"
)

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Download and install talosctl for the current OS/ARCH",
	Long: `Download the talosctl binary for the current OS/ARCH at the specified version.
	
This binary will be saved into the path specified by the "bin-path" flag.
It will be named "talosctl-$version".

Make sure to check the "use <version>" command after installing a new version`,
	Args:  cobra.ExactArgs(1),
	RunE: installVersion,
}

func installVersion(cmd *cobra.Command, args []string) error {
	vrsPath, err := utils.GetVrsPath(cmd)
	if err != nil {
		cmd.Println("Error getting vrs path:", err)
		return err
	}
	err = utils.EnsurePathExists(vrsPath)
	if err != nil {
		cmd.Println("Error ensuring vrs path exists:", err)
		return err
	}
	// TODO check if the version is already installed

	
	ghc := github.New()
	if err := ghc.DownloadRelease(args[0], vrsPath); err != nil {
		return err
	}
	
	cmd.Printf("talosctl version %s successfully installed\n", args[0])
	
	// NOTE: in the following we do not return the error, since it's best effort
	use, err := cmd.Flags().GetBool("use")
	if err != nil {
		cmd.Println("Error retrieving 'use' flag:", err)
		cmd.Println("Skipping action")
		cmd.Printf("To switch to that version run `talosctlenv use %s`\n", args[0])
		return nil
	}
	if !use {
		cmd.Printf("To switch to that version run `talosctlenv use %s`\n", args[0])
		return nil
	}
	// here we want to use the install version immediately
	rootCmd.SetArgs([]string{"use", args[0]})
	if err := rootCmd.Execute(); err != nil {
		cmd.Println("Error executing use action:", err)
		cmd.Println("Skipping action")
		cmd.Printf("To switch to that version run `talosctlenv use %s`\n", args[0])
		return nil
	}

	return nil
}

func init() {
	installCmd.LocalFlags().BoolP("use", "u", false, "Immediately use the version once installed (best effort)")
	rootCmd.AddCommand(installCmd)
}
