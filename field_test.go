package core_test

import (
	"testing"

	"github.com/newlix/core"
	"github.com/tj/assert"
)

func TestBuiltinTypeFields(t *testing.T) {
	t.Run("filters to builtin types only", func(t *testing.T) {
		custom := mustInitTypes(t, core.Type{Name: "address"})[0]
		fields := []core.Field{
			{Name: "id", Type: core.Int},
			{Name: "name", Type: core.String},
			{Name: "address", Type: custom},
			{Name: "active", Type: core.Bool},
		}
		got := core.BuiltinTypeFields(fields)
		assert.Equal(t, 3, len(got))
		assert.Equal(t, "id", got[0].Name)
		assert.Equal(t, "name", got[1].Name)
		assert.Equal(t, "active", got[2].Name)
	})

	t.Run("returns empty for no builtin fields", func(t *testing.T) {
		custom := mustInitTypes(t, core.Type{Name: "tag"})[0]
		fields := []core.Field{
			{Name: "tag", Type: custom},
		}
		got := core.BuiltinTypeFields(fields)
		assert.Equal(t, 0, len(got))
	})

	t.Run("returns empty for nil input", func(t *testing.T) {
		got := core.BuiltinTypeFields(nil)
		assert.Equal(t, 0, len(got))
	})
}
