package step

// TerraformDestroy creates a step that runs terraform destroy with the provided arguments.
func TerraformDestroy(workspace string, args ...string) Step {
	args = append([]string{"destroy", "-auto-approve"}, args...)
	command := createCommand("terraform", args...)
	command.Env = append(command.Env, tfWorkspaceEnv(workspace))
	return runCommand("could not destroy: %v", command)
}
