package step

// ModInit creates a step that initializes the go module.
func ModInit(finalModulePath *string) Step {
	return func() error {
		return runCommand(
			"failed to initialize go module: %v",
			createCommand("go", "mod", "init", *finalModulePath),
		)()
	}
}
