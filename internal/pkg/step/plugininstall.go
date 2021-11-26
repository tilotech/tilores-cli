package step

import (
	"fmt"
	"os"
	"strings"
)

func PluginInstall(pkg, version, target string) Step {
	return func() error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		cmd := createCommand("go", "install", fmt.Sprintf("%v@%v", pkg, version))
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOBIN=%v", wd))

		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to get plugin dependency %v: %v", pkg, err)
		}

		sourceParts := strings.Split(pkg, "/")
		source := sourceParts[len(sourceParts)-1]
		err = os.Rename(source, target)
		if err != nil {
			return fmt.Errorf("failed to move plugin dependency %v from %v to %v: %v", pkg, source, target, err)
		}

		return nil
	}
}
