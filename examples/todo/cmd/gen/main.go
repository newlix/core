package main

import (
	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/golang"
	"github.com/newlix/core/generators/kotlin"
	"github.com/newlix/core/generators/swift"
)

func main() {
	golang.GenerateTypesFile(golang.GenerateTypesFileConfig{
		Output:  "types.gen.go",
		Package: "todo",
		Types:   spec.Types,
	})
	golang.GenerateClientFile(golang.GenerateClientFileConfig{
		Output:  "client/client.gen.go",
		Package: "client",
		Methods: spec.Methods,
		Types:   spec.Types,
	})
	golang.GenerateServerFile(golang.GenerateServerFileConfig{
		Output:  "server/server.gen.go",
		Package: "server",
		Methods: spec.Methods,
		Types:   spec.Types,
	})

	swift.GenerateTypesFile(swift.GenerateTypesFileConfig{
		Output: "swift/types.gen.swift",
		Types:  spec.Types,
	})

	swift.GenerateClientFile(swift.GenerateClientFileConfig{
		Output:  "swift/client.gen.swift",
		Methods: spec.Methods,
		Types:   spec.Types,
		Client:  "TodoClient",
	})

	kotlin.GenerateTypesFile(kotlin.GenerateTypesFileConfig{
		Output:  "kotlin/types.gen.kt",
		Package: "com.example",
		Types:   spec.Types,
	})

	kotlin.GenerateClientFile(kotlin.GenerateClientFileConfig{
		Output:       "kotlin/client.gen.kt",
		Package:      "com.example",
		Methods:      spec.Methods,
		TypesPackage: "",
		Types:        spec.Types,
		Client:       "TodoClient",
	})
}
