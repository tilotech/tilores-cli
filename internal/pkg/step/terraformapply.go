package step

func TerraformApply(workspace string, args ...string) Step {
	args = append([]string{"apply", "-auto-approve"}, args...)
	command := createCommand("terraform", args...)
	command.Env = append(command.Env, "TF_WORKSPACE="+workspace)
	return runCommand("could not apply: %v", command)
}
