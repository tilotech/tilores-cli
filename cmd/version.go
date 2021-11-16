package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		err := printVersion()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("could not read BuildInfo")
	}
	fmt.Println(buildInfo.Main.Version)
	return nil
}
