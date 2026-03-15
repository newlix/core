package core_test

import (
	"testing"

	"github.com/newlix/core"
	"github.com/tj/assert"
)

func mustInitMethods(t *testing.T, mm ...core.Method) []core.Method {
	t.Helper()
	result, err := core.InitMethods(mm...)
	assert.NoError(t, err)
	return result
}

func TestInitMethods(t *testing.T) {
	t.Run("sorts by name", func(t *testing.T) {
		mm := mustInitMethods(t,
			core.Method{Name: "remove_item"},
			core.Method{Name: "add_item"},
		)
		assert.Equal(t, "add_item", mm[0].Name)
		assert.Equal(t, "remove_item", mm[1].Name)
	})

	t.Run("sets camel names", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{Name: "get_items"})
		assert.Equal(t, "GetItems", mm[0].CamelName)
		assert.Equal(t, "getItems", mm[0].LowerCamelName)
	})

	t.Run("preserves explicit camel names", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name:           "get_items",
			CamelName:      "FetchItems",
			LowerCamelName: "fetchItems",
		})
		assert.Equal(t, "FetchItems", mm[0].CamelName)
		assert.Equal(t, "fetchItems", mm[0].LowerCamelName)
	})

	t.Run("initializes input field types", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name: "add_item",
			Inputs: []core.Field{
				{Name: "item_name", Type: core.String},
			},
		})
		assert.Equal(t, "ItemName", mm[0].Inputs[0].CamelName)
		assert.Equal(t, "itemName", mm[0].Inputs[0].LowerCamelName)
		assert.True(t, mm[0].Inputs[0].Type.IsInitialized())
	})

	t.Run("initializes output field types", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name: "get_count",
			Outputs: []core.Field{
				{Name: "total_count", Type: core.Int},
			},
		})
		assert.Equal(t, "TotalCount", mm[0].Outputs[0].CamelName)
		assert.Equal(t, "totalCount", mm[0].Outputs[0].LowerCamelName)
		assert.True(t, mm[0].Outputs[0].Type.IsInitialized())
	})

	t.Run("initializes method with no inputs or outputs", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name: "health_check",
		})
		assert.Equal(t, 0, len(mm[0].Inputs))
		assert.Equal(t, 0, len(mm[0].Outputs))
		assert.Equal(t, "HealthCheck", mm[0].CamelName)
	})

	t.Run("initializes method with all builtin types", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name: "create_product",
			Inputs: []core.Field{
				{Name: "name", Type: core.String},
				{Name: "price", Type: core.Float},
				{Name: "is_active", Type: core.Bool},
				{Name: "count", Type: core.Int},
			},
			Outputs: []core.Field{
				{Name: "id", Type: core.Int},
			},
		})
		assert.Equal(t, 4, len(mm[0].Inputs))
		assert.Equal(t, 1, len(mm[0].Outputs))
		assert.Equal(t, "float64", mm[0].Inputs[1].Type.GoType)
		assert.Equal(t, "bool", mm[0].Inputs[2].Type.GoType)
	})

	t.Run("initializes method with array outputs", func(t *testing.T) {
		mm := mustInitMethods(t, core.Method{
			Name: "list_items",
			Outputs: []core.Field{
				{Name: "items", Type: core.String, IsArray: true},
			},
		})
		assert.True(t, mm[0].Outputs[0].IsArray)
	})

	t.Run("returns error on empty name", func(t *testing.T) {
		_, err := core.InitMethods(core.Method{Name: ""})
		assert.EqualError(t, err, "method name must not be empty")
	})

	t.Run("returns error on duplicate names", func(t *testing.T) {
		_, err := core.InitMethods(
			core.Method{Name: "add_item"},
			core.Method{Name: "add_item"},
		)
		assert.Error(t, err)
	})

	t.Run("returns error on empty input field name", func(t *testing.T) {
		_, err := core.InitMethods(core.Method{
			Name: "add_item",
			Inputs: []core.Field{
				{Name: "", Type: core.String},
			},
		})
		assert.Error(t, err)
	})

	t.Run("returns error on empty output field name", func(t *testing.T) {
		_, err := core.InitMethods(core.Method{
			Name: "add_item",
			Outputs: []core.Field{
				{Name: "", Type: core.Int},
			},
		})
		assert.Error(t, err)
	})
}
