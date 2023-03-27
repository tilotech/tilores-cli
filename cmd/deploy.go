package cmd

import (
	"fmt"
	"os"

	"github.com/tilotech/tilores-cli/internal/pkg/step"

	"github.com/spf13/cobra"
)

var (
	region    string
	profile   string
	workspace string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys " + applicationName + " in your own AWS account.",
	Long:  "Deploys " + applicationName + " in your own AWS account.",
	Run: func(cmd *cobra.Command, args []string) {
		err := deployTiloRes(true)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = deployCmd.MarkPersistentFlagRequired("region")

	deployCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")

	deployCmd.PersistentFlags().StringVar(&workspace, "workspace", "default", "The deployments workspace/environment e.g. dev, prod.")
}

func deployTiloRes(apply bool) error {
	dir, err := os.MkdirTemp("", applicationNameLower)
	if err != nil {
		return err
	}

	deployArgs := []string{
		"-var", fmt.Sprintf("profile=%s", profile),
		"-var", fmt.Sprintf("region=%s", region),
		"-var", fmt.Sprintf("api_file=%s/api.zip", dir),
		"-var", fmt.Sprintf("rule_config_file=%s/rule-config.zip", dir),
	}
	var deployStep step.Step
	if apply {
		deployStep = step.TerraformApply(workspace, deployArgs...)
	} else {
		deployStep = step.TerraformPlan(workspace, deployArgs...)
	}

	steps := []step.Step{
		step.TerraformRequire,
		step.Generate,
		step.Build("./cmd/api/...", dir+"/api", "GOOS=linux", "GOARCH=arm64"),
		step.PackageZipRename(dir+"/api", dir+"/api.zip", "bootstrap"),
		step.PackageZip("./rule-config.json", dir+"/rule-config.zip"),
		step.Chdir("deployment/tilores"),
		step.TerraformInit,
		step.TerraformNewWorkspace(workspace),
		deployStep,
	}

	return step.Execute(steps)
}
