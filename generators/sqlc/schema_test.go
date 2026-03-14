package sqlc_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/sqlc"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateSchema(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_schema.gen.sql")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var buf bytes.Buffer
	sqlc.GenerateSchema(&buf, spec.Types)
	got := buf.String()
	w, err := os.Create("testdata/todo_schema.gen.sql")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
