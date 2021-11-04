package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	dir, err := os.MkdirTemp("", "tilores-init")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cases := map[string]struct {
		args      []string
		expectErr bool
	}{
		"create project": {
			args:      []string{dir},
			expectErr: false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {

			err := initializeProject(nil, c.args)

			if c.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
