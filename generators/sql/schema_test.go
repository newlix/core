package sql_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/sql"
)

func TestGenerateSchema(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	sql.GenerateSchema(&act, spec.Types, sql.Cockroachdb)
	got := act.String()
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
