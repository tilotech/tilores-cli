package cmd

import (
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func createGoCommand(args ...string) *exec.Cmd {
	command := exec.Command("go", args...) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func randLowerCaseString(n int) string {
	// TODO: remove the next line, when making the next breaking change and we're requiring at least go v1.20
	rand.Seed(time.Now().UnixNano()) //nolint:staticcheck
	b := make([]byte, n)
	l := int64(len(letters))
	for i := range b {
		b[i] = letters[rand.Int63()%l] //nolint:gosec
	}
	return string(b)
}
