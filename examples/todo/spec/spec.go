package spec

import (
	"log"

	"github.com/newlix/core"
	"github.com/newlix/core/examples/todo/spec/methods"
	"github.com/newlix/core/examples/todo/spec/types"
)

var (
	Types   []core.Type
	Methods []core.Method
)

func init() {
	var err error
	Types, err = core.InitTypes(
		types.Item,
	)
	if err != nil {
		log.Fatal(err)
	}
	Methods, err = core.InitMethods(
		methods.AddItem,
		methods.GetItems,
		methods.RemoveItem,
	)
	if err != nil {
		log.Fatal(err)
	}
}
