package step

func TerraformApply(args ...string) Step {
	args = append([]string{"apply", "-auto-approve"}, args...)
	return runCommand(
		"could not apply: %v",
		"terraform", args...,
	)
}
