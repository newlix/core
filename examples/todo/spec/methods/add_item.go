package methods

import (
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec/types"
)

var AddItem = core.Method{
	Name:        "add_item",
	Description: "AddItem adds an item to the list.",
	Inputs: []core.Field{
		{
			Name:        "item",
			Description: "the item to add.",
			Type:        types.Item,
		},
	},
}
