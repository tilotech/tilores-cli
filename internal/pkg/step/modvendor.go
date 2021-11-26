package step

func ModVendor() error {
	return runCommand(
		"failed to vendor project dependencies: %v",
		"go", "mod", "vendor",
	)()
}
