package step

// TerraformRequire creates a step that ensures that terraform is installed.
func TerraformRequire() error {
	return runCommand(
		"terraform is required, but was not found: %v",
		createCommand("terraform", "version"),
	)()
}
