package step

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

// PackageZipRename packages source file into target zip, the resulting zip contains the source file
// but named with newSourceName
func PackageZipRename(source string, target string, newSourceName string) Step {
	return func() error {
		sourceFile, err := os.Open(source) //nolint:gosec // reason: static path
		if err != nil {
			return err
		}
		defer sourceFile.Close() //nolint:gosec,errcheck // reason: opened for read

		targetFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer targetFile.Close() //nolint:gosec,errcheck // reason: only applicable for error cases

		zipWriter := zip.NewWriter(targetFile)
		zipFile, err := zipWriter.Create(path.Base(newSourceName))
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, sourceFile)
		if err != nil {
			return err
		}
		err = zipWriter.Close()
		if err != nil {
			return err
		}

		if err = targetFile.Close(); err != nil {
			return err
		}

		return nil
	}
}
