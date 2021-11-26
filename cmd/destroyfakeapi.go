package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
)

// destroyFakeApiCmd represents the fakeApi command
var destroyFakeApiCmd = &cobra.Command{
	Use:   "fake-api",
	Short: "Removes the previously deployed fake API from your AWS account.",
	Long:  "Removes the previously deployed fake API from your AWS account.",
	Run: func(cmd *cobra.Command, args []string) {
		err := destroyFakeAPI()
		cobra.CheckErr(err)
	},
}

func init() {
	destroyCmd.AddCommand(destroyFakeApiCmd)
}

func destroyFakeAPI() error {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	steps := []step.Step{
		step.TerraformRequire,
		step.Chdir("deployment/fake-api"),

		// For some reason Terraform requires the variables being set during destroy.
		// See: https://github.com/hashicorp/terraform/issues/23552
		//
		// Additionally the lambda module checks also on destroy if the files exists.
		// Therefore we must provide an empty file as input.
		step.TerraformDestroy(
			"-var", fmt.Sprintf("profile=%s", profile),
			"-var", fmt.Sprintf("region=%v", region),
			"-var", fmt.Sprintf("api_file=%v", f.Name()),
			"-var", fmt.Sprintf("dispatcher_file=%v", f.Name()),
		),
	}

	return step.Execute(steps)
}
