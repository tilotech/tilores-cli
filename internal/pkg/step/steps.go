package step

type Step func() error

func Execute(steps []Step) error {
	for _, f := range steps {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
