package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tilotech/tilores-cli/internal/pkg/step"

	"github.com/spf13/cobra"
)

var (
	region    string
	profile   string
	workspace string
	varFile   string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys " + applicationName + " in your own AWS account.",
	Long:  "Deploys " + applicationName + " in your own AWS account.",
	Run: func(_ *cobra.Command, _ []string) {
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

	deployCmd.PersistentFlags().StringVar(&varFile, "var-file", "", "The path to the file that holds the values for terraform variables")
}

func deployTiloRes(apply bool) error {
	dir, err := os.MkdirTemp("", applicationNameLower)
	if err != nil {
		return err
	}

	externalRefs, err := extractExternalRefLists()
	if err != nil {
		return fmt.Errorf("failed to extract external reference lists: %w", err)
	}
	refsJSON, err := json.Marshal(externalRefs)
	if err != nil {
		return fmt.Errorf("failed to marshal external reference lists: %w", err)
	}

	deployArgs := []string{
		"-var", fmt.Sprintf("profile=%s", profile),
		"-var", fmt.Sprintf("region=%s", region),
		"-var", fmt.Sprintf("api_file=%s/api.zip", dir),
		"-var", fmt.Sprintf("rule_config_file=%s/rule-config.zip", dir),
		"-var", fmt.Sprintf("external_reflists=%s", refsJSON),
	}
	if varFile != "" {
		deployArgs = append(deployArgs, fmt.Sprintf("-var-file=%s", varFile))
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

func extractExternalRefLists() ([]string, error) {
	data, err := os.ReadFile("./rule-config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // File doesn't exist, no external refs
		}
		return nil, err
	}

	var config struct {
		ReferenceLists []struct {
			External string `json:"external"`
		} `json:"referenceLists"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var refs []string
	for _, ref := range config.ReferenceLists {
		if ref.External != "" {
			refs = append(refs, ref.External)
		}
	}
	return refs, nil
}
