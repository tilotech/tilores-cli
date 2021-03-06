package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/tilotech/tilores-cli/internal/pkg/step"
	"github.com/tilotech/tilores-cli/templates"
)

func CreateUpgradeSteps(upgradeVersion string, variables map[string]interface{}) ([]step.Step, error) {
	stepFiles, err := templates.Upgrades.ReadDir(fmt.Sprintf("upgrades/%v", upgradeVersion))
	if err != nil {
		return nil, err
	}
	sort.Sort(ByName(stepFiles))

	steps := []step.Step{}
	for _, stepFile := range stepFiles {
		if stepFile.IsDir() {
			continue
		}
		fn := stepFile.Name()
		parts := strings.SplitN(fn, "_", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("internal upgrade error: invalid file %v", fn)
		}
		action := parts[1]
		stepFile := fmt.Sprintf("upgrades/%v/%v", upgradeVersion, fn)
		switch {
		case strings.HasPrefix(action, "dependencies"):
			addDepsSteps, err := createAddDependencySteps(stepFile)
			if err != nil {
				return nil, err
			}
			steps = append(steps, addDepsSteps...)
		case strings.HasPrefix(action, "replace"):
			replaceSteps, err := createReplaceSteps(stepFile, variables)
			if err != nil {
				return nil, err
			}
			steps = append(steps, replaceSteps...)
		case strings.HasPrefix(action, "install_plugin"):
			installPluginSteps, err := createInstallPluginSteps(stepFile)
			if err != nil {
				return nil, err
			}
			steps = append(steps, installPluginSteps...)
		default:
			return nil, fmt.Errorf("internal upgrade error: invalid file %v", fn)
		}
	}

	return steps, nil
}

type ByName []fs.DirEntry

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func createAddDependencySteps(fileName string) ([]step.Step, error) {
	deps := []string{}
	err := decodeStepFile(fileName, &deps)
	if err != nil {
		return nil, err
	}

	return []step.Step{
		step.GetDependencies(deps),
	}, nil
}

func createReplaceSteps(fileName string, variables map[string]interface{}) ([]step.Step, error) {
	replace := &struct {
		Target      string
		OldTemplate string
		NewTemplate string
	}{}

	err := decodeStepFile(fileName, replace)
	if err != nil {
		return nil, err
	}

	data, err := templates.Upgrades.ReadFile(replace.OldTemplate)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("t").Parse(string(data))
	if err != nil {
		return nil, err
	}

	expectedFileContent := &bytes.Buffer{}
	err = tmpl.Execute(expectedFileContent, variables)
	if err != nil {
		return nil, err
	}

	targetFile, err := os.Open(replace.Target)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = targetFile.Close()
	}()
	actualFileContent, err := ioutil.ReadAll(targetFile)
	if err != nil {
		return nil, err
	}

	steps := []step.Step{}
	if !bytes.Equal(expectedFileContent.Bytes(), actualFileContent) {
		steps = append(
			steps,
			step.Backup(replace.Target),
		)
	}

	steps = append(
		steps,
		step.RenderTemplate(templates.Upgrades, replace.NewTemplate, replace.Target, variables),
	)

	return steps, nil
}

func createInstallPluginSteps(fileName string) ([]step.Step, error) {
	install := &struct {
		Pkg     string
		Version string
		Target  string
	}{}

	err := decodeStepFile(fileName, install)
	if err != nil {
		return nil, err
	}

	return []step.Step{
		step.PluginInstall(install.Pkg, install.Version, install.Target),
	}, nil
}

func decodeStepFile(fileName string, v interface{}) error {
	f, err := templates.Upgrades.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	decoder := json.NewDecoder(f)

	return decoder.Decode(v)
}
