package step

import "fmt"

// TerraformApply creates a step that runs terraform apply with the provided arguments.
func TerraformApply(workspace string, args ...string) Step {
	args = append([]string{"apply", "-auto-approve"}, args...)
	command := createCommand("terraform", args...)
	command.Env = append(command.Env, tfWorkspaceEnv(workspace))
	return runCommand("could not apply: %v", command)
}

func tfWorkspaceEnv(workspace string) string {
	return fmt.Sprintf("TF_WORKSPACE=%v", workspace)
}
