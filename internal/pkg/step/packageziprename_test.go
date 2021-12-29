package step

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageZipRename(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = ioutil.WriteFile(dir+"/somefile", []byte("some data"), 0666)
	require.NoError(t, err)

	newSourceName := "newSourceName"
	fixture := PackageZipRename(dir+"/somefile", dir+"/somefile.zip", newSourceName)
	actual := fixture()

	readerCloser, err := zip.OpenReader(dir + "/somefile.zip")
	require.NoError(t, err)
	assert.Equal(t, newSourceName, readerCloser.File[0].Name)
	defer readerCloser.Close()

	assert.NoError(t, actual)
	assert.FileExists(t, dir+"/somefile.zip")
}

func TestPackageZipRenameWithNonExistingSourceFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	fixture := PackageZipRename(dir+"/somefile", dir+"/somefile.zip", "somefile")
	actual := fixture()

	assert.Error(t, actual)
	assert.NoFileExists(t, dir+"/somefile.zip")
}

func TestPackageZipRenameWithInvalidTargetFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	err = ioutil.WriteFile(dir+"/somefile", []byte("some data"), 0666)
	require.NoError(t, err)

	fixture := PackageZipRename(dir+"/somefile", "/dev/null/is/not/a/valid/file/name", "somefile")
	actual := fixture()

	assert.Error(t, actual)
}
