/*
Copyright © 2021 Tilo Tech GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/tilotech/tilores-cli/templates"
)

var (
	modulePath        string
	dispatcherVersion string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new " + applicationName + " application",
	Long: `Initalize (` + applicationNameLower + ` init) will create a new ` + applicationName + ` application and
the appropriate structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeProject(args)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&modulePath, "module-path", "m", "", "The go module path for the generated go.mod file, defaults to the project folder name")
	initCmd.Flags().StringVar(&dispatcherVersion, "dispatcher-version", "latest", "The version of the fake dispatcher plugin used for local runs")
}

func initializeProject(args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	err := os.MkdirAll(path, 0750)
	if err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	err = os.Chdir(path)
	if err != nil {
		return err
	}

	finalModulePath := modulePath
	if finalModulePath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		finalModulePath = filepath.Base(wd)
	}

	err = createGoCommand("mod", "init", finalModulePath).Run()
	if err != nil {
		return fmt.Errorf("failed to initialize go module: %v", err)
	}

	variables := templateVariables{
		ApplicationName: applicationName,
		GeneratedMsg:    generatedMsg,
		ModulePath:      finalModulePath,
	}

	err = copyTemplatesRecursive(templates.InitPreGenerate, "", variables)
	if err != nil {
		return err
	}

	err = getDependencies([]string{
		"gitlab.com/tilotech/tilores-plugin-api/dispatcher",
	})
	if err != nil {
		return err
	}

	err = createGoCommand("mod", "vendor").Run()
	if err != nil {
		return fmt.Errorf("failed to vendor project dependencies: %v", err)
	}

	err = createGoCommand("generate", "./...").Run()
	if err != nil {
		return fmt.Errorf("failed to generate project resources: %v", err)
	}

	err = copyTemplatesRecursive(templates.InitPostGenerate, "", variables)
	if err != nil {
		return err
	}

	err = getPluginDependencies()
	if err != nil {
		return err
	}

	err = createGoCommand("mod", "tidy").Run()
	if err != nil {
		return fmt.Errorf("failed to vendor project dependencies: %v", err)
	}

	err = createGoCommand("mod", "vendor").Run()
	if err != nil {
		return fmt.Errorf("failed to vendor project dependencies: %v", err)
	}

	err = createGoCommand("build").Run()
	if err != nil {
		return fmt.Errorf("failed to verify project by running go build: %v", err)
	}

	return nil
}

func copyTemplatesRecursive(fs embed.FS, path string, variables templateVariables) error {
	templateFiles, err := fs.ReadDir("init" + path)
	if err != nil {
		return err
	}
	for _, file := range templateFiles {
		filePath := fmt.Sprintf("%s/%s", path, file.Name())
		if file.IsDir() {
			err = os.MkdirAll("."+filePath, 0750)
			if err != nil {
				return err
			}
			err = copyTemplatesRecursive(fs, filePath, variables)
			if err != nil {
				return err
			}
		} else {
			err = copyTemplateFile(fs, filePath, variables)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyTemplateFile(fs embed.FS, path string, variables templateVariables) error {
	data, err := fs.ReadFile("init" + path)
	if err != nil {
		return err
	}

	tmpl, err := template.New("t").Parse(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse template file %v: %v", path, err)
	}

	errCh := make(chan error, 1)
	r, w := io.Pipe()
	go func() {
		defer close(errCh)
		defer func() {
			_ = w.Close()
		}()
		err := tmpl.Execute(w, variables)
		if err != nil {
			errCh <- err
			return
		}
	}()

	data, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}

	return os.WriteFile("."+path[:len(path)-len(".tmpl")], data, 0600)
}

func getDependencies(dependencies []string) error {
	for _, dependency := range dependencies {
		err := createGoCommand("get", dependency).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func getPluginDependencies() error {
	err := getPluginDependency("gitlab.com/tilotech/tilores-plugin-fake-dispatcher", dispatcherVersion, "tilores-plugin-dispatcher")
	if err != nil {
		return err
	}

	return nil
}

func getPluginDependency(pkg, version, target string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get plugin dependency %v: %v", pkg, err)
	}

	cmd := createGoCommand("install", fmt.Sprintf("%v@%v", pkg, version))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOBIN=%v", wd))

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to get plugin dependency %v: %v", pkg, err)
	}

	sourceParts := strings.Split(pkg, "/")
	source := sourceParts[len(sourceParts)-1]
	err = os.Rename(source, target)
	if err != nil {
		return fmt.Errorf("failed to move plugin dependency %v from %v to %v: %v", pkg, source, target, err)
	}

	return nil
}

type templateVariables struct {
	ApplicationName string
	GeneratedMsg    string
	ModulePath      string
}
