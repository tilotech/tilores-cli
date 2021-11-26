package step

func BuildVerify() error {
	return runCommand(
		"failed to verify project by running go build: %v",
		"go", "build", "./...",
	)()
}
