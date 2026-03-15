package common

import (
	"io"
	"os"
	"path"
)

// errWriter wraps an io.Writer and captures the first write error.
// Subsequent writes after an error are no-ops.
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) Write(p []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	n, err := ew.w.Write(p)
	ew.err = err
	return n, err
}

// GenerateFile creates the output file (and parent directories) then calls fn to write content.
// Write errors from fn are captured even if fn does not check them (e.g. via common.Out).
func GenerateFile(output string, fn func(w io.Writer) error) error {
	if err := os.MkdirAll(path.Dir(output), 0o700); err != nil {
		return err
	}
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	ew := &errWriter{w: f}
	fnErr := fn(ew)
	if fnErr != nil || ew.err != nil {
		f.Close()
		os.Remove(output)
		if fnErr != nil {
			return fnErr
		}
		return ew.err
	}
	return f.Close()
}
