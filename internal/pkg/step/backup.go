package step

import (
	"fmt"
	"os"
)

const colorReset = "\033[0m"
const colorYellow = "\033[33m"

// Backup creates a step that backups the provided file path.
func Backup(path string) Step {
	return func() error {
		fmt.Printf("%vDetected a modified file. Please compare the differences between %v and %v.bak.%v\n", colorYellow, path, path, colorReset)
		return os.Rename(path, fmt.Sprintf("%v.bak", path))
	}
}
