package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades to the latest TiloRes version",
	Long:  "Upgrades to the latest TiloRes version.",
	Run: func(cmd *cobra.Command, args []string) {
		err := upgrade()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() error {
	finalModulePath, err := pkg.GetModulePath()
	if err != nil {
		return err
	}
	variables := map[string]interface{}{
		"ApplicationName": applicationName,
		"GeneratedMsg":    generatedMsg,
		"ModulePath":      finalModulePath,
	}
	config, err := pkg.LoadConfig()
	if err != nil {
		return err
	}
	upgrades, err := pkg.ListUpgrades(config.Version)
	if err != nil {
		return err
	}
	if len(upgrades) == 0 {
		fmt.Println("Upgrade not required. Already at the latest version.")
		return nil
	}
	for _, upgradeVersion := range upgrades {
		steps, err := pkg.CreateUpgradeSteps(upgradeVersion, variables)
		if err != nil {
			return fmt.Errorf("error while preparing the upgrade for %v: %v", upgradeVersion, err)
		}
		steps = append(
			steps,
			step.ModTidy,
			step.ModVendor,
			step.Generate,
			step.BuildVerify,
		)
		err = step.Execute(steps)
		if err != nil {
			return fmt.Errorf("error while running the upgrade for %v: %v", upgradeVersion, err)
		}
		config.Version = upgradeVersion
		err = pkg.SaveConfig(config)
		if err != nil {
			return err
		}
	}
	fmt.Println("Successfully upgraded to the latest version.")
	return nil
}
