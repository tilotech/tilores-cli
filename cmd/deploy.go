package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	region  string
	profile string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys " + applicationName + " in your own AWS account.",
	Long: `Deploys ` + applicationName + ` in your own AWS account.

By default it deploys the full application. Using "deploy fake-api" you can
deploy only the API with a fake implementation, similar to the run command.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called")
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = deployCmd.MarkPersistentFlagRequired("region")

	deployCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")
}
