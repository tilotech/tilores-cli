package step

func TerraformInit() error {
	return runCommand(
		"could not initialize terraform: %v",
		"terraform", "init",
	)()
}
