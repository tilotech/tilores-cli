package step

import "io"

// TerraformNewWorkspace creates a step that initializes a new terraform workspace if it does not yet exist.
func TerraformNewWorkspace(workspace string) Step {
	return func() error {
		command := createCommand("terraform", "workspace", "new", workspace)
		command.Stderr = io.Discard
		_ = command.Run() // the error is ignored because the command is expected to fail if the workspace already exists
		return nil
	}
}
