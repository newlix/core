package core_test

import (
	"testing"

	"github.com/newlix/core"
	"github.com/tj/assert"
)

func mustInitTypes(t *testing.T, tt ...core.Type) []core.Type {
	t.Helper()
	result, err := core.InitTypes(tt...)
	assert.NoError(t, err)
	return result
}

func TestInitTypes(t *testing.T) {
	t.Run("sorts by name", func(t *testing.T) {
		tt := mustInitTypes(t,
			core.Type{Name: "zebra"},
			core.Type{Name: "apple"},
		)
		assert.Equal(t, "apple", tt[0].Name)
		assert.Equal(t, "zebra", tt[1].Name)
	})

	t.Run("sets camel names", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{Name: "user_profile"})
		assert.Equal(t, "UserProfile", tt[0].CamelName)
		assert.Equal(t, "userProfile", tt[0].LowerCamelName)
	})

	t.Run("preserves explicit camel names", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name:           "item",
			CamelName:      "MyItem",
			LowerCamelName: "myItem",
		})
		assert.Equal(t, "MyItem", tt[0].CamelName)
		assert.Equal(t, "myItem", tt[0].LowerCamelName)
	})

	t.Run("defaults primary key to id", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{Name: "item"})
		assert.Equal(t, []string{"id"}, tt[0].PrimaryKey)
	})

	t.Run("preserves explicit primary key", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name:       "item",
			PrimaryKey: []string{"slug"},
		})
		assert.Equal(t, []string{"slug"}, tt[0].PrimaryKey)
	})

	t.Run("marks as initialized", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{Name: "item"})
		assert.True(t, tt[0].IsInitialized())
	})

	t.Run("initializes field types recursively", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name: "order",
			Fields: []core.Field{
				{Name: "total", Type: core.Float},
			},
		})
		assert.True(t, tt[0].Fields[0].Type.IsInitialized())
	})

	t.Run("sets field camel names", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name: "item",
			Fields: []core.Field{
				{Name: "created_at", Type: core.String},
			},
		})
		assert.Equal(t, "CreatedAt", tt[0].Fields[0].CamelName)
		assert.Equal(t, "createdAt", tt[0].Fields[0].LowerCamelName)
	})

	t.Run("builtin types report as builtin", func(t *testing.T) {
		assert.True(t, core.String.IsBuiltin())
		assert.True(t, core.Int.IsBuiltin())
		assert.True(t, core.Bool.IsBuiltin())
		assert.True(t, core.Float.IsBuiltin())
	})

	t.Run("custom types are not builtin", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{Name: "item"})
		assert.Equal(t, false, tt[0].IsBuiltin())
	})

	t.Run("initializes type with all builtin field types", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name: "product",
			Fields: []core.Field{
				{Name: "id", Type: core.Int},
				{Name: "name", Type: core.String},
				{Name: "is_active", Type: core.Bool},
				{Name: "price", Type: core.Float},
			},
		})
		assert.Equal(t, 4, len(tt[0].Fields))
		assert.Equal(t, "int", tt[0].Fields[0].Type.GoType)
		assert.Equal(t, "string", tt[0].Fields[1].Type.GoType)
		assert.Equal(t, "bool", tt[0].Fields[2].Type.GoType)
		assert.Equal(t, "float64", tt[0].Fields[3].Type.GoType)
	})

	t.Run("initializes type with no fields", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{Name: "empty"})
		assert.Equal(t, 0, len(tt[0].Fields))
		assert.True(t, tt[0].IsInitialized())
	})

	t.Run("preserves array field flag", func(t *testing.T) {
		tt := mustInitTypes(t, core.Type{
			Name: "collection",
			Fields: []core.Field{
				{Name: "tags", Type: core.String, IsArray: true},
			},
		})
		assert.True(t, tt[0].Fields[0].IsArray)
	})

	t.Run("initializes custom type field", func(t *testing.T) {
		address := core.Type{
			Name:   "address",
			Fields: []core.Field{{Name: "street", Type: core.String}},
		}
		tt := mustInitTypes(t,
			address,
			core.Type{
				Name: "person",
				Fields: []core.Field{
					{Name: "home", Type: address},
				},
			},
		)
		person := tt[1] // sorted: address, person
		assert.Equal(t, "Address", person.Fields[0].Type.CamelName)
		assert.True(t, person.Fields[0].Type.IsInitialized())
		assert.False(t, person.Fields[0].Type.IsBuiltin())
	})

	t.Run("returns error on empty name", func(t *testing.T) {
		_, err := core.InitTypes(core.Type{Name: ""})
		assert.EqualError(t, err, "type name must not be empty")
	})

	t.Run("returns error on duplicate names", func(t *testing.T) {
		_, err := core.InitTypes(
			core.Type{Name: "item"},
			core.Type{Name: "item"},
		)
		assert.Error(t, err)
	})

	t.Run("returns error on empty field name", func(t *testing.T) {
		_, err := core.InitTypes(core.Type{
			Name: "item",
			Fields: []core.Field{
				{Name: "", Type: core.String},
			},
		})
		assert.Error(t, err)
	})
}
