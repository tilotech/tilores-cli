package step

func Generate() error {
	return runCommand(
		"failed to generate project resources: %v",
		"go", "generate", "./...",
	)()
}
