package kotlin_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/kotlin"
)

func TestGenerateClient(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_client.kt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	kotlin.GenerateClient(&act, spec.Methods, "TodoClient")
	got := act.String()
	w, err := os.Create("testdata/todo_client.gen.kt")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
