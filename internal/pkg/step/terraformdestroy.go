package step

func TerraformDestroy(args ...string) Step {
	args = append([]string{"destroy", "-auto-approve"}, args...)
	return runCommand(
		"could not destroy: %v",
		"terraform", args...,
	)
}
