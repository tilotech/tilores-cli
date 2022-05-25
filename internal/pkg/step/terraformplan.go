package step

func TerraformPlan(workspace string, args ...string) Step {
	args = append([]string{"plan"}, args...)
	command := createCommand("terraform", args...)
	command.Env = append(command.Env, "TF_WORKSPACE="+workspace)
	return runCommand("could not plan: %v", command)
}
