package cmd

import (
	"os"
	"os/exec"
)

func createGoCommand(args ...string) *exec.Cmd {
	command := exec.Command("go", args...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command
}

type step func() error

func executeSteps(steps []step) error {
	for _, f := range steps {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
