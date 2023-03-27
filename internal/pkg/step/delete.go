package step

import "os"

// Delete creates a step that removes the provided file.
func Delete(path string) Step {
	return func() error {
		return os.Remove(path)
	}
}
