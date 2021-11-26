package step

import "os"

func Chdir(path string) Step {
	return func() error {
		return os.Chdir(path)
	}
}
