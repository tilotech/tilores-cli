package step

import "os"

// Chdir creates a step that changes the directory to the provided path.
func Chdir(path string) Step {
	return func() error {
		return os.Chdir(path)
	}
}
