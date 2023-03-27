package pkg

import (
	"encoding/json"
	"os/exec"
)

// GetModulePath returns the current module path of the project.
func GetModulePath() (string, error) {
	cmd := exec.Command("go", "mod", "edit", "-json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	mp := struct {
		Module struct {
			Path string
		}
	}{}
	err = json.Unmarshal(out, &mp)
	if err != nil {
		return "", err
	}
	return mp.Module.Path, nil
}
