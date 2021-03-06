package step

func ModInit(finalModulePath *string) Step {
	return func() error {
		return runCommand(
			"failed to initialize go module: %v",
			createCommand("go", "mod", "init", *finalModulePath),
		)()
	}
}
