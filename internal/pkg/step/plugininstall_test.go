package step

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.Chdir(dir)
	require.NoError(t, err)

	cases := map[string]struct {
		pkg         string
		version     string
		target      string
		expectFile  string
		expectError bool
	}{
		"install plugin": {
			pkg:        "github.com/tilotech/tilores-plugin-fake-dispatcher",
			version:    "latest",
			target:     "dispatcher",
			expectFile: dir + "/dispatcher",
		},
		"invalid package": {
			pkg:         "this is not a valid package",
			version:     "latest",
			target:      "dispatcher",
			expectError: true,
		},
		"invalid target": {
			pkg:         "github.com/tilotech/tilores-plugin-fake-dispatcher",
			version:     "latest",
			target:      "..",
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			err := PluginInstall(c.pkg, c.version, c.target)()

			if c.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.FileExists(t, c.expectFile)
			}
		})
	}
}
