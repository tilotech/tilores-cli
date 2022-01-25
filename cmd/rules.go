package cmd

import (
	"github.com/spf13/cobra"
)

// rulesSimulateCmd represents simulateRules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "A base for all rules related commands.",
	Long:  `A base for all rules related commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(rulesCmd)
}
