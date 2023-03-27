package step

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitProjectPath creates a step that gathers information about the project path.
func InitProjectPath(path string, finalModulePath *string) Step {
	return func() error {
		err := os.MkdirAll(path, 0750)
		if err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}

		err = os.Chdir(path)
		if err != nil {
			return err
		}

		if *finalModulePath == "" {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %v", err)
			}
			*finalModulePath = filepath.Base(wd)
		}
		return nil
	}
}
