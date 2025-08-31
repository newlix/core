package kotlin_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/kotlin"
	"github.com/google/go-cmp/cmp"
)

// func TestGenerateTypes(t *testing.T) {
// 	b, err := os.ReadFile("testdata/todo_types.kt")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	want := string(b)
// 	var act bytes.Buffer
// 	if err := kotlin.GenerateTypes(&act, spec.Types); err != nil {
// 		t.Error(err)
// 	}
// 	got := act.String()

// 	if got != want {
// 		t.Error(cmp.Diff(got, want))
// 		fmt.Println(got)
// 	}
// }

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
