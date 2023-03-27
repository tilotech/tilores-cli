package step

// TerraformInit creates a step that runs terraform init with the upgrade flag.
func TerraformInit() error {
	command := createCommand("terraform", "init", "-upgrade")
	command.Env = append(command.Env, "TF_WORKSPACE=default")
	return runCommand("could not initialize terraform: %v", command)()
}

// TerraformInitFast creates a step that runs terraform init without the upgrade flag.
func TerraformInitFast() error {
	command := createCommand("terraform", "init")
	command.Env = append(command.Env, "TF_WORKSPACE=default")
	return runCommand("could not initialize terraform: %v", command)()
}
