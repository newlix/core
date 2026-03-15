package common

import (
	"io"
	"os"
	"path"
)

// GenerateFile creates the output file (and parent directories) then calls fn to write content.
func GenerateFile(output string, fn func(w io.Writer) error) error {
	if err := os.MkdirAll(path.Dir(output), 0o700); err != nil {
		return err
	}
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := fn(w); err != nil {
		os.Remove(output)
		return err
	}
	return nil
}
