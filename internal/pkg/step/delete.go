package step

import "os"

func Delete(path string) Step {
	return func() error {
		return os.Remove(path)
	}
}
