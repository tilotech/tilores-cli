package step

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

// PackageZip creates a step that packages the source into a zip file.
func PackageZip(source string, target string) Step {
	return func() error {
		sourceFile, err := os.Open(source) //nolint:gosec // reason: static path
		if err != nil {
			return err
		}
		defer sourceFile.Close() //nolint:gosec,errcheck // reason: opened for read

		targetFile, err := os.Create(target) //nolint:gosec
		if err != nil {
			return err
		}
		defer targetFile.Close() //nolint:gosec,errcheck // reason: only applicable for error cases

		zipWriter := zip.NewWriter(targetFile)
		zipFile, err := zipWriter.Create(path.Base(source))
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

		return targetFile.Close()
	}
}
