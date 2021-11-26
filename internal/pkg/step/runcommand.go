package step

import (
	"fmt"
	"os"
	"os/exec"
)

func runCommand(errMsg string, name string, args ...string) Step {
	return func() error {
		command := createCommand(name, args...)
		err := command.Run()
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}
		return nil
	}
}

func createCommand(name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command
}
