package step

import "fmt"

func Build(pkg string, target string) Step {
	return runCommand(
		fmt.Sprintf("could not build %v: %%v", pkg),
		"go", "build", "-o", target, pkg,
	)
}
