package step

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitProjectPath(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cases := map[string]struct {
		path               string
		finalModulePath    string
		expectedModulePath string
		expectError        bool
	}{
		"empty module path": {
			path:               dir + "/empty-module-path",
			finalModulePath:    "",
			expectedModulePath: "empty-module-path",
		},
		"path with slash": {
			path:               dir + "/path-with-slash/",
			finalModulePath:    "",
			expectedModulePath: "path-with-slash",
		},
		"with module path": {
			path:               dir + "/with-module-path",
			finalModulePath:    "some-module-path",
			expectedModulePath: "some-module-path",
		},
		"invalid path": {
			path:        "/dev/null/invalid/path",
			expectError: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := c.finalModulePath
			err := InitProjectPath(c.path, &actual)()

			if c.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.expectedModulePath, actual)
			}
		})
	}
}
