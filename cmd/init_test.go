package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	cases := map[string]struct {
		args                 []string
		modulePathFlag       string
		expectFilesExist     []string
		expectFilesToContain map[string]string
	}{
		"create project with path argument": {
			args: []string{"foobar"},
			expectFilesExist: []string{
				"foobar/go.mod",
				"foobar/foobar",
				"foobar/gqlgen.yml",
				"foobar/server.go",
			},
			expectFilesToContain: map[string]string{
				"foobar/go.mod": "module foobar",
			},
		},
		"create project without path argument with module path": {
			args:           []string{},
			modulePathFlag: "example.com/test/foopkg",
			expectFilesExist: []string{
				"go.mod",
				"foopkg",
				"gqlgen.yml",
				"server.go",
			},
			expectFilesToContain: map[string]string{
				"go.mod": "module example.com/test/foopkg",
			},
		},
		"create project without path argument or flags": {
			args: []string{},
			expectFilesExist: []string{
				"go.mod",
				"gqlgen.yml",
				"server.go",
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dir, err := createTempDir()
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			modulePath = c.modulePathFlag

			err = initializeProject(c.args)
			assert.NoError(t, err)

			for _, file := range c.expectFilesExist {
				assert.FileExists(t, dir+"/"+file)
			}

			for file, expectedPartialContent := range c.expectFilesToContain {
				actualContent, err := ioutil.ReadFile(dir + "/" + file)
				require.NoError(t, err)
				assert.Contains(t, string(actualContent), expectedPartialContent)
			}
		})
	}
}

func createTempDir() (string, error) {
	dir, err := os.MkdirTemp("", applicationNameLower)
	if err != nil {
		return "", err
	}

	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}
