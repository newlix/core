package common

import (
	"fmt"
	"io"
)

// Out writes a formatted line to w.
func Out(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, format+"\n", a...)
}
