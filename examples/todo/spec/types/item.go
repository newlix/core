package types

import "github.com/newlix/core"

var Item = core.Type{
	Name:        "item",
	Description: "Item is a to-do item.",
	Fields: []core.Field{
		{
			Name:        "id",
			CamelName:   "ID",
			Description: "ID is the unique id",
			Type:        core.Int,
		},
		{
			Name:        "text",
			Description: "Text is the content",
			Type:        core.String,
		},
		{
			Name:        "created_at",
			Description: "CreatedAt is the timestamp which the item created at",
			Type:        core.Time,
		},
	},
	GoPackage: "github.com/newlix/core/examples/todo",
}
