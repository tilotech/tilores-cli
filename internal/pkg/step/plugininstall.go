package step

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

		// Install plugin for linux os (used by AWS lambda)
		cmdLinux := createCommand("go", "install", fmt.Sprintf("%v@%v", pkg, version))
		linuxGoEnvs := []string{fmt.Sprintf("GOPATH=%v/plugins", wd), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0"}
		cmdLinux.Env = os.Environ()
		cmdLinux.Env = append(cmdLinux.Env, linuxGoEnvs...)

		err = cmdLinux.Run()
		if err != nil {
			return fmt.Errorf("failed to get plugin dependency for linux %v: %v", pkg, err)
		}

		pluginLinuxBinPath := wd + "/plugins/bin/linux_amd64/" + source
		if _, err := os.Stat(pluginLinuxBinPath); errors.Is(err, os.ErrNotExist) {
			pluginLinuxBinPath = wd + "/plugins/bin/" + source
		}
		err = os.Rename(pluginLinuxBinPath, target+"-linux-amd64")
		if err != nil {
			return fmt.Errorf("failed to move plugin dependency %v from %v to %v: %v", pkg, pluginLinuxBinPath, target+"-linux-amd64", err)
		}

		_, err = exec.Command("chmod", "-R", "0755", wd+"/plugins").CombinedOutput() //nolint:gosec
		if err != nil {
			return fmt.Errorf("chmod failed on plugins folder: %v", err)
		}
		err = os.RemoveAll(wd + "/plugins/")
		if err != nil {
			return fmt.Errorf("failed to cleanup after plugins install: %v", err)
		}

		return nil
	}
}
