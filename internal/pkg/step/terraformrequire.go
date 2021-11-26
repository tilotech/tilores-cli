package step

func TerraformRequire() error {
	return runCommand(
		"terraform is required, but was not found: %v",
		"terraform", "version",
	)()
}
