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
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	modulePath string

	//go:embed templates/schema templates/tools templates/generate.go.tmpl templates/gqlgen.yml.tmpl
	templatesPreGenerate embed.FS

	//go:embed templates/server.go.tmpl
	templatesPostGenerate embed.FS
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
}

func initializeProject(args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	err := os.MkdirAll(path, 0755)
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

	out, err := exec.Command("go", "mod", "init", finalModulePath).CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return fmt.Errorf("failed to initialize go module: %v", err)
	}

	variables := templateVariables{
		ApplicationName: applicationName,
		GeneratedMsg:    generatedMsg,
		ModulePath:      finalModulePath,
	}

	err = copyTemplatesRecursive(templatesPreGenerate, "", variables)
	if err != nil {
		return err
	}

	out, err = exec.Command("go", "mod", "vendor").CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return fmt.Errorf("failed to vendor project dependencies: %v", err)
	}

	out, err = exec.Command("go", "generate", "./...").CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return fmt.Errorf("failed to generate project resources: %v", err)
	}

	err = copyTemplatesRecursive(templatesPostGenerate, "", variables)
	if err != nil {
		return err
	}

	out, err = exec.Command("go", "build").CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return fmt.Errorf("failed to verify project by running go build: %v", err)
	}

	return nil
}

func copyTemplatesRecursive(fs embed.FS, path string, variables templateVariables) error {
	templateFiles, err := fs.ReadDir("templates" + path)
	if err != nil {
		return err
	}
	for _, file := range templateFiles {
		filePath := fmt.Sprintf("%s/%s", path, file.Name())
		if file.IsDir() {
			err = os.MkdirAll("."+filePath, 0755)
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
	data, err := fs.ReadFile("templates" + path)
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
		defer w.Close()
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

	return os.WriteFile("."+path[:len(path)-len(".tmpl")], data, 0644)
}

type templateVariables struct {
	ApplicationName string
	GeneratedMsg    string
	ModulePath      string
}
