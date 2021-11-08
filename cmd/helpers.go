package cmd

import (
	"os"
	"os/exec"
)

func createGoCommand(args ...string) *exec.Cmd {
	command := exec.Command("go", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command
}
