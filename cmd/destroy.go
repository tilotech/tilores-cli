package cmd

import (
	"fmt"
	"os"

	"github.com/tilotech/tilores-cli/internal/pkg/step"

	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Removes " + applicationName + " from your AWS account.",
	Long:  "Removes " + applicationName + " from your AWS account.",
	Run: func(_ *cobra.Command, _ []string) {
		err := destroyTiloRes()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = destroyCmd.MarkPersistentFlagRequired("region")

	destroyCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")

	destroyCmd.PersistentFlags().StringVar(&workspace, "workspace", "default", "The deployments workspace/environment e.g. dev, prod.")

	destroyCmd.PersistentFlags().StringVar(&varFile, "var-file", "", "The path to the file that holds the values for terraform variables")
}

func destroyTiloRes() error {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	destroyArgs := []string{
		"-var", fmt.Sprintf("profile=%s", profile),
		"-var", fmt.Sprintf("region=%v", region),
		"-var", fmt.Sprintf("api_file=%v", f.Name()),
		"-var", fmt.Sprintf("rule_config_file=%v", f.Name()),
	}
	if varFile != "" {
		destroyArgs = append(destroyArgs, fmt.Sprintf("-var-file=%s", varFile))
	}

	steps := []step.Step{
		step.TerraformRequire,
		step.Chdir("deployment/tilores"),
		step.TerraformInitFast,
		step.TerraformNewWorkspace(workspace),

		// For some reason Terraform requires the variables being set during destroy.
		// See: https://github.com/hashicorp/terraform/issues/23552
		//
		// Additionally the lambda module checks also on destroy if the files exists.
		// Therefore we must provide an empty file as input.
		step.TerraformDestroy(workspace, destroyArgs...),
	}

	return step.Execute(steps)
}
