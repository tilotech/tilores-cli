package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
)

// deployFakeApiCmd represents the fake-api command
var deployFakeApiCmd = &cobra.Command{
	Use:   "fake-api",
	Short: "Deploys only the API into your AWS account together with a fake implementation.",
	Long:  "Deploys only the API into your AWS account together with a fake implementation.",
	Run: func(cmd *cobra.Command, args []string) {
		err := deployFakeAPI()
		cobra.CheckErr(err)
	},
}

func init() {
	deployCmd.AddCommand(deployFakeApiCmd)
}

func deployFakeAPI() error {
	dir, err := os.MkdirTemp("", applicationNameLower)
	if err != nil {
		return err
	}

	steps := []step.Step{
		step.TerraformRequire,
		step.Generate,
		step.Build("./cmd/api/...", dir+"/api", "GOOS=linux", "GOARCH=arm64"),
		step.PackageZipRename(dir+"/api", dir+"/api.zip", "bootstrap"),
		step.PackageZipRename("tilores-plugin-fake-dispatcher-linux-arm64", dir+"/dispatcher.zip", "tilores-plugin-fake-dispatcher"),
		step.Chdir("deployment/fake-api"),
		step.TerraformInit,
		step.TerraformApply(
			"default",
			"-var", fmt.Sprintf("profile=%s", profile),
			"-var", fmt.Sprintf("region=%s", region),
			"-var", fmt.Sprintf("api_file=%s/api.zip", dir),
			"-var", fmt.Sprintf("dispatcher_file=%s/dispatcher.zip", dir),
		),
	}

	return step.Execute(steps)
}
