package step

func TerraformInit() error {
	command := createCommand("terraform", "init", "-upgrade")
	command.Env = append(command.Env, "TF_WORKSPACE=default")
	return runCommand("could not initialize terraform: %v", command)()
}
