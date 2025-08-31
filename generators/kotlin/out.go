package kotlin

import (
	"fmt"
	"io"
)

func out(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, format+"\n", a...)
}
