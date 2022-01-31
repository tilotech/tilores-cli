package step

import "os"

func TerraformInit() error {
	command := createCommand("terraform", "init")
	command.Env = os.Environ()
	command.Env = append(command.Env, "TF_WORKSPACE=default")
	return runCommand("could not initialize terraform: %v", command)()
}
