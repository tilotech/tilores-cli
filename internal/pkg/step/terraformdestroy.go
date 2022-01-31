package step

import "os"

func TerraformDestroy(workspace string, args ...string) Step {
	args = append([]string{"destroy", "-auto-approve"}, args...)
	command := createCommand("terraform", args...)
	command.Env = os.Environ()
	command.Env = append(command.Env, "TF_WORKSPACE="+workspace)
	return runCommand("could not destroy: %v", command)
}
