package golang_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/golang"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateTypes(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_types.go")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	golang.GenerateTypes(&act, "github.com/newlix/core/examples/todo", spec.Types, []string{"json", "db"})
	got := act.String()
	w, err := os.Create("testdata/todo_types.gen.go")
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
	b, err := os.ReadFile("testdata/todo_method_types.go")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	golang.GenerateMethodTypes(&act, "github.com/newlix/core/examples/todo/client", spec.Methods, spec.Types)
	got := act.String()
	w, err := os.Create("testdata/todo_method_types.gen.go")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
