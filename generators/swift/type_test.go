package swift_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/swift"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateTypes(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_types.swift")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	swift.GenerateTypes(&act, spec.Types)
	got := act.String()
	w, err := os.Create("testdata/todo_types.gen.swift")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateMethodTypes(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_method_types.swift")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	swift.GenerateMethodTypes(&act, spec.Methods)
	got := act.String()
	w, err := os.Create("testdata/todo_method_types.gen.swift")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
