package cmd

import (
	"fmt"
	"os"

	"github.com/tilotech/tilores-cli/internal/pkg/step"

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
		err := deployTiloRes()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVar(&region, "region", "", "The deployments AWS region.")
	_ = deployCmd.MarkPersistentFlagRequired("region")

	deployCmd.PersistentFlags().StringVar(&profile, "profile", "", "The AWS credentials profile.")
}

func deployTiloRes() error {
	dir, err := os.MkdirTemp("", applicationNameLower)
	if err != nil {
		return err
	}

	steps := []step.Step{
		step.TerraformRequire,
		step.Generate,
		step.Build("./cmd/api/...", dir+"/api", "GOOS=linux", "GOARCH=arm64"),
		step.PackageZipRename(dir+"/api", dir+"/api.zip", "bootstrap"),
		step.PackageZip("./rule-config.json", dir+"/rule-config.zip"),
		step.Chdir("deployment/tilores"),
		step.TerraformInit,
		step.TerraformApply(
			"-var", fmt.Sprintf("profile=%s", profile),
			"-var", fmt.Sprintf("region=%s", region),
			"-var", fmt.Sprintf("api_file=%s/api.zip", dir),
			"-var", fmt.Sprintf("rule_config_file=%s/rule-config.zip", dir),
		),
	}

	return step.Execute(steps)
}
