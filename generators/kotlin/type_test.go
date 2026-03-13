package kotlin_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/kotlin"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateTypes(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_types.kt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	kotlin.GenerateTypes(&act, spec.Types)
	got := act.String()
	w, err := os.Create("testdata/todo_types.gen.kt")
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
	b, err := os.ReadFile("testdata/todo_method_types.kt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	kotlin.GenerateMethodTypes(&act, spec.Methods, spec.Types)
	got := act.String()
	w, err := os.Create("testdata/todo_method_types.gen.kt")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateTypes_EmptyFields(t *testing.T) {
	types := core.InitTypes(core.Type{
		Name:        "empty",
		Description: "empty type",
	})
	var buf bytes.Buffer
	kotlin.GenerateTypes(&buf, types)
	got := buf.String()
	want := "// empty type\n@Serializable\nclass Empty(\n)\n\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateMethodTypes_EmptyOutputs(t *testing.T) {
	methods := core.InitMethods(core.Method{
		Name:        "ping",
		Description: "health check",
	})
	var buf bytes.Buffer
	kotlin.GenerateMethodTypes(&buf, methods, nil)
	got := buf.String()
	want := "@Serializable\nclass PingInput(\n)\n\n@Serializable\nclass PingOutput(\n)\n\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
