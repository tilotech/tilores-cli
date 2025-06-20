package pkg

import (
	"fmt"
	"os"
	"regexp"
)

type ProjectVariables struct {
	DeployPrefix string
}

var deployPrefixRegex = regexp.MustCompile(`resource_prefix\s*=\s*\"([^\"]+)\"`)

func CollectVariables() (*ProjectVariables, error) {
	mainTFFile, err := os.ReadFile("deployment/tilores/main.tf")
	if err != nil {
		return nil, err
	}

	matches := deployPrefixRegex.FindStringSubmatch((string(mainTFFile)))
	if len(matches) != 2 {
		return nil, fmt.Errorf("could not match regex to find resource_prefix in main.tf")
	}
	return &ProjectVariables{
		DeployPrefix: matches[1],
	}, nil
}
