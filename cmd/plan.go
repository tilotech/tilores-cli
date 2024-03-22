package cmd

import (
	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plans the " + applicationName + " deployment for your own AWS account.",
	Long:  `Plans the ` + applicationName + ` deployment for your own AWS account.`,
	Run: func(_ *cobra.Command, _ []string) {
		err := deployTiloRes(false)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = planCmd.MarkPersistentFlagRequired("region")

	planCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")

	planCmd.PersistentFlags().StringVar(&workspace, "workspace", "default", "The deployments workspace/environment e.g. dev, prod.")

	planCmd.PersistentFlags().StringVar(&tfvarsFile, "tfvars-file", "", "The path to the file that holds the values for terraform variables")
}
