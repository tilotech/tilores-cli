package step

func GetDependencies(dependencies []string) func() error {
	return func() error {
		for _, dependency := range dependencies {
			err := runCommand(
				"failed to get dependencies: %v",
				"go", "get", dependency,
			)()
			if err != nil {
				return err
			}
		}
		return nil
	}
}
