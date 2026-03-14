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

func mustInitTypes(t *testing.T, tt ...core.Type) []core.Type {
	t.Helper()
	result, err := core.InitTypes(tt...)
	assert.NoError(t, err)
	return result
}

func mustInitMethods(t *testing.T, mm ...core.Method) []core.Method {
	t.Helper()
	result, err := core.InitMethods(mm...)
	assert.NoError(t, err)
	return result
}

func TestGoName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Id", "ID"},
		{"UserId", "UserID"},
		{"HttpUrl", "HTTPURL"},
		{"ApiKey", "APIKey"},
		{"Item", "Item"},
		{"AddItem", "AddItem"},
		{"GetItems", "GetItems"},
		{"CreatedAt", "CreatedAt"},
		{"JsonRpc", "JSONRPC"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, golang.GoName(tt.input))
		})
	}
}

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
	golang.GenerateMethodTypes(&act, "github.com/newlix/core/examples/todo/client", spec.Methods)
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

func TestGenerateTypes_AllBuiltinTypes(t *testing.T) {
	types := mustInitTypes(t, core.Type{
		Name:        "product",
		Description: "Product is an item for sale.",
		Fields: []core.Field{
			{Name: "id", Description: "unique identifier", Type: core.Int},
			{Name: "name", Description: "product name", Type: core.String},
			{Name: "is_active", Description: "whether the product is active", Type: core.Bool},
			{Name: "price", Description: "product price", Type: core.Float},
		},
		GoPackage: "example.com/shop",
	})
	var buf bytes.Buffer
	golang.GenerateTypes(&buf, "example.com/shop", types, []string{"json", "db"})
	got := buf.String()

	assert.Contains(t, got, "type Product struct {")
	assert.Contains(t, got, "ID int")
	assert.Contains(t, got, "Name string")
	assert.Contains(t, got, "IsActive bool")
	assert.Contains(t, got, "Price float64")
	assert.Contains(t, got, `json:"id"`)
	assert.Contains(t, got, `json:"is_active"`)
	assert.Contains(t, got, `json:"price"`)
}

func TestGenerateTypes_EmptyFields(t *testing.T) {
	types := mustInitTypes(t, core.Type{
		Name:        "empty",
		Description: "Empty has no fields.",
		GoPackage:   "example.com/shop",
	})
	var buf bytes.Buffer
	golang.GenerateTypes(&buf, "example.com/shop", types, []string{"json", "db"})
	got := buf.String()
	want := "// Empty has no fields.\ntype Empty struct {\n}\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateTypes_ArrayFields(t *testing.T) {
	types := mustInitTypes(t, core.Type{
		Name:        "order",
		Description: "Order contains items.",
		Fields: []core.Field{
			{Name: "id", Description: "order id", Type: core.Int},
			{Name: "tags", Description: "order tags", Type: core.String, IsArray: true},
		},
		GoPackage: "example.com/shop",
	})
	var buf bytes.Buffer
	golang.GenerateTypes(&buf, "example.com/shop", types, []string{"json", "db"})
	got := buf.String()

	assert.Contains(t, got, "Tags []string")
	assert.Contains(t, got, "ID int")
}

func TestGenerateTypes_CustomTypeField(t *testing.T) {
	address := core.Type{
		Name:        "address",
		Description: "Address is a postal address.",
		Fields: []core.Field{
			{Name: "street", Description: "street name", Type: core.String},
		},
		GoPackage: "example.com/shop",
	}
	types := mustInitTypes(t,
		address,
		core.Type{
			Name:        "person",
			Description: "Person is a user.",
			Fields: []core.Field{
				{Name: "name", Description: "person name", Type: core.String},
				{Name: "home", Description: "home address", Type: address},
			},
			GoPackage: "example.com/shop",
		},
	)
	var buf bytes.Buffer
	golang.GenerateTypes(&buf, "example.com/shop", types, []string{"json", "db"})
	got := buf.String()

	assert.Contains(t, got, "type Address struct {")
	assert.Contains(t, got, "type Person struct {")
	assert.Contains(t, got, "Home Address")
}

func TestGenerateMethodTypes_NoInputsNoOutputs(t *testing.T) {
	methods := mustInitMethods(t, core.Method{
		Name:        "health_check",
		Description: "HealthCheck returns server status.",
	})
	var buf bytes.Buffer
	golang.GenerateMethodTypes(&buf, "example.com/shop/client", methods)
	got := buf.String()
	want := "type HealthCheckInput struct {\n}\n\ntype HealthCheckOutput struct {\n}\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateMethodTypes_AllBuiltinInputs(t *testing.T) {
	methods := mustInitMethods(t, core.Method{
		Name:        "create_product",
		Description: "CreateProduct creates a product.",
		Inputs: []core.Field{
			{Name: "name", Description: "product name", Type: core.String},
			{Name: "price", Description: "product price", Type: core.Float},
			{Name: "is_active", Description: "active flag", Type: core.Bool},
			{Name: "count", Description: "item count", Type: core.Int},
		},
	})
	var buf bytes.Buffer
	golang.GenerateMethodTypes(&buf, "example.com/shop/client", methods)
	got := buf.String()

	assert.Contains(t, got, "type CreateProductInput struct {")
	assert.Contains(t, got, "Name string")
	assert.Contains(t, got, "Price float64")
	assert.Contains(t, got, "IsActive bool")
	assert.Contains(t, got, "Count int")
	assert.Contains(t, got, "type CreateProductOutput struct {")
}
