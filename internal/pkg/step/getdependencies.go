package step

func GetDependencies(dependencies []string) func() error {
	return func() error {
		for _, dependency := range dependencies {
			err := runCommand(
				"failed to get dependencies: %v",
				createCommand("go", "get", "-d", dependency),
			)()
			if err != nil {
				return err
			}
		}
		return nil
	}
}
