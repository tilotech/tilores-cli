package step

import (
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

func RenderTemplates(fs embed.FS, prefix string, variables interface{}) Step {
	return func() error {
		return copyTemplatesRecursive(fs, "", prefix, variables)
	}
}

func copyTemplatesRecursive(fs embed.FS, path string, prefix string, variables interface{}) error {
	templateFiles, err := fs.ReadDir(prefix + path)
	if err != nil {
		return err
	}
	for _, file := range templateFiles {
		filePath := fmt.Sprintf("%s/%s", path, file.Name())
		if file.IsDir() {
			err = os.MkdirAll("."+filePath, 0750)
			if err != nil {
				return err
			}
			err = copyTemplatesRecursive(fs, filePath, prefix, variables)
			if err != nil {
				return err
			}
		} else {
			err = copyTemplateFile(fs, filePath, prefix, variables)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyTemplateFile(fs embed.FS, path string, prefix string, variables interface{}) error {
	data, err := fs.ReadFile(prefix + path)
	if err != nil {
		return err
	}

	tmpl, err := template.New("t").Parse(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse template file %v: %v", path, err)
	}

	errCh := make(chan error, 1)
	r, w := io.Pipe()
	go func() {
		defer close(errCh)
		defer func() {
			_ = w.Close()
		}()
		err := tmpl.Execute(w, variables)
		if err != nil {
			errCh <- err
			return
		}
	}()

	data, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if err, ok := <-errCh; ok {
		return err
	}

	return os.WriteFile("."+path[:len(path)-len(".tmpl")], data, 0600)
}
