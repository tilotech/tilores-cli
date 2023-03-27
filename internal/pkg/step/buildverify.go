package step

// BuildVerify creates a step that ensures the validity of the generated project.
func BuildVerify() error {
	return runCommand(
		"failed to verify project by running go build: %v",
		createCommand("go", "build", "./..."),
	)()
}
