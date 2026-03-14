package swift_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/swift"
	"github.com/tj/assert"
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

func TestGenerateTypes_AllBuiltinTypes(t *testing.T) {
	types, err := core.InitTypes(core.Type{
		Name:        "product",
		Description: "Product is an item for sale.",
		Fields: []core.Field{
			{Name: "id", Description: "unique identifier", Type: core.Int},
			{Name: "name", Description: "product name", Type: core.String},
			{Name: "is_active", Description: "active flag", Type: core.Bool},
			{Name: "price", Description: "product price", Type: core.Float},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	swift.GenerateTypes(&buf, types)
	got := buf.String()

	assert.Contains(t, got, "struct Product: Codable {")
	assert.Contains(t, got, "var id: Int = 0")
	assert.Contains(t, got, `var name: String = ""`)
	assert.Contains(t, got, "var isActive: Bool = false")
	assert.Contains(t, got, "var price: Double = 0.0")
	assert.Contains(t, got, `case isActive = "is_active"`)
	assert.Contains(t, got, "extension Product {")
	assert.Contains(t, got, "decodeIfPresent(Bool.self, forKey: .isActive)")
	assert.Contains(t, got, "decodeIfPresent(Double.self, forKey: .price)")
}

func TestGenerateTypes_EmptyFields(t *testing.T) {
	types, err := core.InitTypes(core.Type{
		Name:        "empty",
		Description: "Empty has no fields.",
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	swift.GenerateTypes(&buf, types)
	got := buf.String()
	want := "// Empty has no fields.\nstruct Empty: Codable {\n}\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateTypes_ArrayFields(t *testing.T) {
	types, err := core.InitTypes(core.Type{
		Name:        "order",
		Description: "Order contains items.",
		Fields: []core.Field{
			{Name: "id", Description: "order id", Type: core.Int},
			{Name: "tags", Description: "order tags", Type: core.String, IsArray: true},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	swift.GenerateTypes(&buf, types)
	got := buf.String()

	assert.Contains(t, got, "var tags: [String] = []")
	assert.Contains(t, got, "decodeIfPresent([String].self, forKey: .tags)")
}

func TestGenerateMethodTypes_NoInputsNoOutputs(t *testing.T) {
	methods, err := core.InitMethods(core.Method{
		Name:        "health_check",
		Description: "HealthCheck returns server status.",
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	swift.GenerateMethodTypes(&buf, methods)
	got := buf.String()

	assert.Contains(t, got, "struct HealthCheckInput: Codable {")
	assert.Contains(t, got, "struct HealthCheckOutput: Codable {")
	assert.NotContains(t, got, "extension HealthCheckInput")
	assert.NotContains(t, got, "extension HealthCheckOutput")
}
