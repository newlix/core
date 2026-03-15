package methods

import (
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec/types"
)

var RemoveItem = core.Method{
	Name:        "remove_item",
	Description: "RemoveItem removes an item from the to-do list.",
	Table:       "item",
	Inputs: []core.Field{
		{
			Name:        "id",
			Description: "the id of the item to remove.",
			Type:        core.Int,
		},
	},
	Outputs: []core.Field{
		{
			Name:        "item",
			Description: "the item removed.",
			Type:        types.Item,
		},
	},
}
