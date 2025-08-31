package spec

import (
	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec/methods"
	"github.com/newlix/core/examples/todo/spec/types"
)

var (
	Types   []core.Type
	Methods []core.Method
)

func init() {
	Types = core.InitTypes(
		types.Item,
	)
	Methods = core.InitMethods(
		methods.AddItem,
		methods.GetItems,
		methods.RemoveItem,
	)
}
