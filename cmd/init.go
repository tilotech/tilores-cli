package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
	"github.com/tilotech/tilores-cli/templates"
)

var (
	modulePath        string
	dispatcherVersion string
	deployPrefix      string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initializes a new " + applicationName + " application",
	Long: `Initalize (` + toolName + ` init) will create a new ` + applicationName + ` application and
the appropriate structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeProject(args)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&modulePath, "module-path", "m", "", "The go module path for the generated go.mod file, defaults to the project folder name.")
	initCmd.Flags().StringVar(&dispatcherVersion, "dispatcher-version", "latest", "The version of the fake dispatcher plugin used for local runs.")
	initCmd.Flags().StringVar(&deployPrefix, "deploy-prefix", "", "The initial prefix for resources created during the deploy phase, defaults to a random eight character string and can be changed later in the generated files.")
}

func initializeProject(args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	finalModulePath := modulePath

	finalDeployPrefix := deployPrefix
	if finalDeployPrefix == "" {
		finalDeployPrefix = randLowerCaseString(8)
	}

	variables := templateVariables{
		ApplicationName: applicationName,
		GeneratedMsg:    generatedMsg,
		ModulePath:      &finalModulePath,
		DeployPrefix:    finalDeployPrefix,
	}

	steps := []step.Step{
		step.InitProjectPath(path, &finalModulePath),
		step.ModInit(&finalModulePath),
		step.RenderTemplates(templates.InitPreGenerate, "init", variables),
		step.GetDependencies([]string{
			"github.com/tilotech/tilores-plugin-api/dispatcher",
		}),
		step.ModVendor,
		step.Generate,
		step.RenderTemplates(templates.InitPostGenerate, "init", variables),
		step.PluginInstall("github.com/tilotech/tilores-plugin-fake-dispatcher", dispatcherVersion, "tilores-plugin-dispatcher"),
		step.ModTidy,
		step.ModVendor,
		step.BuildVerify,
	}

	return step.Execute(steps)
}

type templateVariables struct {
	ApplicationName string
	GeneratedMsg    string
	ModulePath      *string
	DeployPrefix    string
}
