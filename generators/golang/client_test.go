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

func TestGenerateClient(t *testing.T) {
	b, err := os.ReadFile("testdata/todo_client.go")
	if err != nil {
		t.Fatal(err)
	}
	want := string(b)
	var act bytes.Buffer
	golang.GenerateClient(&act, spec.Methods)
	got := act.String()
	w, err := os.Create("testdata/todo_client.gen.go")
	if err != nil {
		t.Fatal(err)
	}
	w.WriteString(got)
	w.Close()
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateClient_NoInputsNoOutputs(t *testing.T) {
	methods := mustInitMethods(t, core.Method{
		Name:        "health_check",
		Description: "HealthCheck returns server status.",
	})
	var buf bytes.Buffer
	golang.GenerateClient(&buf, methods)
	got := buf.String()

	assert.Contains(t, got, "func (c *Client) HealthCheck(in HealthCheckInput) (HealthCheckOutput, error)")
	assert.Contains(t, got, "var out HealthCheckOutput")
	assert.Contains(t, got, `call(c.HTTPClient, c.AuthToken, c.URL, "health_check", in, &out)`)
}

func TestGenerateClient_MultipleMethods(t *testing.T) {
	methods := mustInitMethods(t,
		core.Method{
			Name:        "create_product",
			Description: "CreateProduct creates a product.",
			Inputs: []core.Field{
				{Name: "name", Description: "product name", Type: core.String},
			},
		},
		core.Method{
			Name:        "health_check",
			Description: "HealthCheck returns server status.",
		},
	)
	var buf bytes.Buffer
	golang.GenerateClient(&buf, methods)
	got := buf.String()

	assert.Contains(t, got, "func (c *Client) CreateProduct(")
	assert.Contains(t, got, "func (c *Client) HealthCheck(")
}
