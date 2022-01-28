package step

func Generate() error {
	return runCommand(
		"failed to generate project resources: %v",
		createCommand("go", "generate", "./..."),
	)()
}
