package step

// Generate creates a step that runs go generate for all files.
func Generate() error {
	return runCommand(
		"failed to generate project resources: %v",
		createCommand("go", "generate", "./..."),
	)()
}
