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
	kotlin.GenerateMethodTypes(&act, spec.Methods)
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
	types, err := core.InitTypes(core.Type{
		Name:        "empty",
		Description: "empty type",
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	kotlin.GenerateTypes(&buf, types)
	got := buf.String()
	want := "// empty type\n@Serializable\nclass Empty(\n)\n\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateMethodTypes_EmptyOutputs(t *testing.T) {
	methods, err := core.InitMethods(core.Method{
		Name:        "ping",
		Description: "health check",
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	kotlin.GenerateMethodTypes(&buf, methods)
	got := buf.String()
	want := "@Serializable\nclass PingInput(\n)\n\n@Serializable\nclass PingOutput(\n)\n\n"
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
	kotlin.GenerateTypes(&buf, types)
	got := buf.String()
	want := "// Product is an item for sale.\n" +
		"@Serializable\n" +
		"data class Product(\n" +
		"    @SerialName(\"id\") val id: Long = 0,\n" +
		"    @SerialName(\"name\") val name: String = \"\",\n" +
		"    @SerialName(\"is_active\") val isActive: Boolean = false,\n" +
		"    @SerialName(\"price\") val price: Double = 0.0,\n" +
		")\n\n"
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
	kotlin.GenerateTypes(&buf, types)
	got := buf.String()
	want := "// Order contains items.\n" +
		"@Serializable\n" +
		"data class Order(\n" +
		"    @SerialName(\"id\") val id: Long = 0,\n" +
		"    @SerialName(\"tags\") val tags: List<String> = emptyList(),\n" +
		")\n\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}

func TestGenerateMethodTypes_AllBuiltinInputs(t *testing.T) {
	methods, err := core.InitMethods(core.Method{
		Name:        "create_product",
		Description: "CreateProduct creates a product.",
		Inputs: []core.Field{
			{Name: "name", Description: "product name", Type: core.String},
			{Name: "price", Description: "product price", Type: core.Float},
			{Name: "is_active", Description: "active flag", Type: core.Bool},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	kotlin.GenerateMethodTypes(&buf, methods)
	got := buf.String()
	want := "@Serializable\n" +
		"data class CreateProductInput(\n" +
		"    @SerialName(\"name\") val name: String = \"\",\n" +
		"    @SerialName(\"price\") val price: Double = 0.0,\n" +
		"    @SerialName(\"is_active\") val isActive: Boolean = false,\n" +
		")\n\n" +
		"@Serializable\n" +
		"class CreateProductOutput(\n" +
		")\n\n"
	if got != want {
		t.Error(cmp.Diff(got, want))
	}
}
