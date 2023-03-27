package step

import (
	"fmt"
	"os"
)

// Build creates a step that runs go build.
func Build(pkg string, target string, goEnvs ...string) Step {
	return func() error {
		cmd := createCommand("go", "build", "-o", target, pkg)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, goEnvs...)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("could not build %v: %%v", pkg)
		}
		return nil
	}
}
