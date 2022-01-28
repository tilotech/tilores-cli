package step

import (
	"fmt"
	"os"
	"os/exec"
)

func runCommand(errMsg string, command *exec.Cmd) Step {
	return func() error {
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
