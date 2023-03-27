package step

// ModVendor creates a step that runs go mod vendor.
func ModVendor() error {
	return runCommand(
		"failed to vendor project dependencies: %v",
		createCommand("go", "mod", "vendor"),
	)()
}
