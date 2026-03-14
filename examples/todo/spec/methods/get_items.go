package methods

import (
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec/types"
)

var GetItems = core.Method{
	Name:        "get_items",
	Description: "GetItems returns all items in the list.",
	Outputs: []core.Field{
		{
			Name:        "items",
			Description: "Items is the list of to-do items.",
			Type:        types.Item,
			IsArray:     true,
		},
	},
}
