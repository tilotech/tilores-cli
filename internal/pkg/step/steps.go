package step

// Step defines the function interface for an action to be performed.
type Step func() error

// Execute runs the provided steps.
func Execute(steps []Step) error {
	for _, f := range steps {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
