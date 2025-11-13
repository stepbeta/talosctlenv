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

	
	ghc := github.New()
	if err := ghc.DownloadRelease(args[0], binPath); err != nil {
		return err
	}
	
	cmd.Printf("talosctl version %s successfully installed\n", args[0])
	
	// TODO add "--use" flag
	// rootCmd.SetArgs([]string{"use", args[0]})
	// if err := rootCmd.Execute(); err != nil {fmt.Println("Error executing use:", err)}
	cmd.Printf("To switch to that version run `talosctlenv use %s`\n", args[0])

	return nil
}

func init() {
	rootCmd.AddCommand(installCmd)
}
