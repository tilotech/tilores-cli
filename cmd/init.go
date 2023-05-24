package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tilotech/tilores-cli/internal/pkg"
	"github.com/tilotech/tilores-cli/internal/pkg/step"
	"github.com/tilotech/tilores-cli/templates"
)

var (
	modulePath   string
	deployPrefix string
)

const (
	goPluginVersion  = "v0.1.1"
	pluginAPIVersion = "v0.10.0"
	gqlgenVersion    = "v0.17.24"
	insightsVersion  = "610b82049d9c39ea1ef29bdc7c57322754d25650"
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
			"github.com/tilotech/go-plugin@" + goPluginVersion,
			"github.com/tilotech/tilores-plugin-api/dispatcher@" + pluginAPIVersion,
			"github.com/99designs/gqlgen@" + gqlgenVersion,
			"github.com/tilotech/tilores-insights/record@" + insightsVersion,
		}),
		step.ModVendor,
		step.Generate,
		step.RenderTemplates(templates.InitPostGenerate, "init", variables),
		step.ModTidy,
		step.ModVendor,
		step.BuildVerify,
	}

	err := step.Execute(steps)
	if err != nil {
		return err
	}

	version, err := pkg.LatestUpgradeVersion()
	if err != nil {
		return err
	}
	config := pkg.DefaultConfig()
	config.Version = version
	return pkg.SaveConfig(config)
}

type templateVariables struct {
	ApplicationName string
	GeneratedMsg    string
	ModulePath      *string
	DeployPrefix    string
}
