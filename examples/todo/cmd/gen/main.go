package main

import (
	"log"

	"github.com/newlix/core/examples/todo/spec"
	"github.com/newlix/core/generators/golang"
	"github.com/newlix/core/generators/kotlin"
	"github.com/newlix/core/generators/swift"
)

func main() {
	must(golang.GenerateTypesFile(golang.GenerateTypesFileConfig{
		Output:  "types.gen.go",
		Package: "github.com/newlix/core/examples/todo",
		Types:   spec.Types,
	}))
	must(golang.GenerateClientFile(golang.GenerateClientFileConfig{
		Output:  "client/client.gen.go",
		Package: "client",
		Methods: spec.Methods,
		Types:   spec.Types,
	}))
	must(golang.GenerateServerFile(golang.GenerateServerFileConfig{
		Output:  "server/server.gen.go",
		Package: "server",
		Methods: spec.Methods,
		Types:   spec.Types,
	}))

	must(swift.GenerateTypesFile(swift.GenerateTypesFileConfig{
		Output: "swift/types.gen.swift",
		Types:  spec.Types,
	}))

	must(swift.GenerateClientFile(swift.GenerateClientFileConfig{
		Output:  "swift/client.gen.swift",
		Methods: spec.Methods,
		Client:  "TodoClient",
	}))

	must(kotlin.GenerateTypesFile(kotlin.GenerateTypesFileConfig{
		Output:  "kotlin/types.gen.kt",
		Package: "com.example",
		Types:   spec.Types,
	}))

	must(kotlin.GenerateClientFile(kotlin.GenerateClientFileConfig{
		Output:  "kotlin/client.gen.kt",
		Package: "com.example",
		Methods: spec.Methods,
		Client:  "TodoClient",
	}))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
