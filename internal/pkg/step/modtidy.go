package step

func ModTidy() error {
	return runCommand(
		"failed to tidy project dependencies: %v",
		"go", "mod", "tidy",
	)()
}
