package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Removes " + applicationName + " from your AWS account.",
	Long: `Removes ` + applicationName + ` from your AWS account.

By default it removes the full application. Using "destroy fake-api" you can
remove the fake implementation API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy called")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = destroyCmd.MarkPersistentFlagRequired("region")

	destroyCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")
}
