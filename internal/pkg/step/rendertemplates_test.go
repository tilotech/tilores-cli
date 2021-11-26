package step

import (
	"embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/rendertemplates/simple.txt.tmpl
var fsSimple embed.FS

//go:embed testdata/rendertemplates/variable.txt.tmpl
var fsVariable embed.FS

//go:embed testdata/rendertemplates/invalid.txt.tmpl
var fsInvalid embed.FS

func TestRenderTemplates(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.Chdir(dir)
	require.NoError(t, err)

	err = RenderTemplates(fsSimple, "testdata/rendertemplates", nil)()
	assert.NoError(t, err)
	assert.FileExists(t, dir+"/simple.txt")

	err = RenderTemplates(fsSimple, "testdata", nil)()
	assert.NoError(t, err)
	assert.FileExists(t, dir+"/rendertemplates/simple.txt")

	err = RenderTemplates(
		fsVariable,
		"testdata/rendertemplates",
		map[string]interface{}{"Variable": "value"},
	)()
	assert.NoError(t, err)
	assert.FileExists(t, dir+"/variable.txt")

	err = RenderTemplates(
		fsVariable,
		"testdata/rendertemplates",
		map[string]interface{}{"Variable": func() {}},
	)()
	assert.Error(t, err)

	err = RenderTemplates(fsInvalid, "testdata/rendertemplates", nil)()
	assert.Error(t, err)
}
