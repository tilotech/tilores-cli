package step

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageZip(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = ioutil.WriteFile(dir+"/somefile", []byte("some data"), 0666)
	require.NoError(t, err)

	fixture := PackageZip(dir+"/somefile", dir+"/somefile.zip")
	actual := fixture()

	assert.NoError(t, actual)
	assert.FileExists(t, dir+"/somefile.zip")
}

func TestPackageZipWithNonExistingSourceFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	fixture := PackageZip(dir+"/somefile", dir+"/somefile.zip")
	actual := fixture()

	assert.Error(t, actual)
	assert.NoFileExists(t, dir+"/somefile.zip")
}

func TestPackageZipWithInvalidTargetFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/somefile", []byte("some data"), 0666)
	require.NoError(t, err)

	fixture := PackageZip(dir+"/somefile", "/dev/null/is/not/a/valid/file/name")
	actual := fixture()

	assert.Error(t, actual)
}
