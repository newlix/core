package golang_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/golang"
	"github.com/tj/assert"
)

func TestGenerateServer(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_server.go")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	golang.GenerateServer(&act, spec.Methods)
	got := act.String()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateServer_WithInputs(t *testing.T) {
	methods := mustInitMethods(t, core.Method{
		Name:        "create_product",
		Description: "CreateProduct creates a product.",
		Inputs: []core.Field{
			{Name: "name", Description: "product name", Type: core.String},
		},
	})
	var buf bytes.Buffer
	golang.GenerateServer(&buf, methods)
	got := buf.String()

	assert.Contains(t, got, `case "/create_product":`)
	assert.Contains(t, got, "core.ReadRequest(r, &in)")
	assert.Contains(t, got, "s.CreateProduct(ctx, in)")
}

func TestGenerateServer_NoInputs(t *testing.T) {
	methods := mustInitMethods(t, core.Method{
		Name:        "health_check",
		Description: "HealthCheck returns server status.",
	})
	var buf bytes.Buffer
	golang.GenerateServer(&buf, methods)
	got := buf.String()

	assert.Contains(t, got, `case "/health_check":`)
	assert.NotContains(t, got, "core.ReadRequest")
	assert.Contains(t, got, "s.HealthCheck(ctx, in)")
}
