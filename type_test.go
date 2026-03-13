package core_test

import (
	"testing"

	"github.com/newlix/core"
	"github.com/tj/assert"
)

func TestInitTypes(t *testing.T) {
	t.Run("sorts by name", func(t *testing.T) {
		tt := core.InitTypes(
			core.Type{Name: "zebra"},
			core.Type{Name: "apple"},
		)
		assert.Equal(t, "apple", tt[0].Name)
		assert.Equal(t, "zebra", tt[1].Name)
	})

	t.Run("sets camel names", func(t *testing.T) {
		tt := core.InitTypes(core.Type{Name: "user_profile"})
		assert.Equal(t, "UserProfile", tt[0].CamelName)
		assert.Equal(t, "userProfile", tt[0].LowerCamelName)
	})

	t.Run("preserves explicit camel names", func(t *testing.T) {
		tt := core.InitTypes(core.Type{
			Name:           "item",
			CamelName:      "MyItem",
			LowerCamelName: "myItem",
		})
		assert.Equal(t, "MyItem", tt[0].CamelName)
		assert.Equal(t, "myItem", tt[0].LowerCamelName)
	})

	t.Run("defaults primary key to id", func(t *testing.T) {
		tt := core.InitTypes(core.Type{Name: "item"})
		assert.Equal(t, []string{"id"}, tt[0].PrimaryKey)
	})

	t.Run("preserves explicit primary key", func(t *testing.T) {
		tt := core.InitTypes(core.Type{
			Name:       "item",
			PrimaryKey: []string{"slug"},
		})
		assert.Equal(t, []string{"slug"}, tt[0].PrimaryKey)
	})

	t.Run("marks as initialized", func(t *testing.T) {
		tt := core.InitTypes(core.Type{Name: "item"})
		assert.True(t, tt[0].IsInitialized())
	})

	t.Run("initializes field types recursively", func(t *testing.T) {
		tt := core.InitTypes(core.Type{
			Name: "order",
			Fields: []core.Field{
				{Name: "total", Type: core.Float},
			},
		})
		assert.True(t, tt[0].Fields[0].Type.IsInitialized())
	})

	t.Run("sets field camel names", func(t *testing.T) {
		tt := core.InitTypes(core.Type{
			Name: "item",
			Fields: []core.Field{
				{Name: "created_at", Type: core.Time},
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
		assert.True(t, core.Time.IsBuiltin())
	})

	t.Run("custom types are not builtin", func(t *testing.T) {
		tt := core.InitTypes(core.Type{Name: "item"})
		assert.Equal(t, false, tt[0].IsBuiltin())
	})
}
