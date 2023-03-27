package step

// ModTidy creates a step that runs go mod tidy.
func ModTidy() error {
	return runCommand(
		"failed to tidy project dependencies: %v",
		createCommand("go", "mod", "tidy"),
	)()
}
