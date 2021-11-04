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
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	modulePath string

	//go:embed templates
	templates embed.FS
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new TiloRes application",
	Long: `Initalize (tilores init) will create a new TiloRes application and
the appropriate structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initializeProject(cmd, args)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&modulePath, "module-path", "m", "", "The go module path for the generated go.mod file, defaults to the project folder name")
}

func initializeProject(cmd *cobra.Command, args []string) error {
	path := args[0]
	fmt.Println(args)

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
	fmt.Println(string(out))
	if err != nil {
		return fmt.Errorf("failed to initialize go module: %v", err)
	}

	err = copyTemplatesRecursive("")
	if err != nil {
		return err
	}

	out, err = exec.Command("go", "build").CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return fmt.Errorf("failed to verify project by running go build: %v", err)
	}

	return nil
}

func copyTemplatesRecursive(path string) error {
	templateFiles, err := templates.ReadDir("templates" + path)
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
			err = copyTemplatesRecursive(filePath)
			if err != nil {
				return err
			}
		} else {
			err = copyTemplateFile(filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyTemplateFile(path string) error {
	data, err := templates.ReadFile("templates" + path)
	if err != nil {
		return err
	}
	return os.WriteFile("."+path, data, 0644)
}
