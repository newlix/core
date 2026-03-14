package swift

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateTypesFileConfig struct {
	Output string
	Types  []core.Type
}

func GenerateTypesFile(c GenerateTypesFileConfig) {
	if err := os.MkdirAll(path.Dir(c.Output), 0o700); err != nil {
		log.Fatal(err)
	}
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)
	out(w, "import Foundation")
	out(w, "")

	GenerateTypes(w, c.Types)
}

func GenerateTypes(w io.Writer, tt []core.Type) {
	// types
	for _, t := range tt {
		out(w, "// %s", t.Description)
		out(w, "struct %s: Codable {", t.CamelName)
		writeFields(w, t.Fields)
		writeCodingKeys(w, t.Fields)
		out(w, "}")
		out(w, "")
		writeDecoderInit(w, t.CamelName, t.Fields)
	}

}

func GenerateMethodTypes(w io.Writer, ms []core.Method) {
	// methods
	for _, m := range ms {
		out(w, "struct %sInput: Codable {", m.CamelName)
		writeFields(w, m.Inputs)
		out(w, "")
		writeCodingKeys(w, m.Inputs)
		out(w, "}")
		out(w, "")
		writeDecoderInit(w, m.CamelName+"Input", m.Inputs)

		out(w, "")

		out(w, "struct %sOutput: Codable {", m.CamelName)
		writeFields(w, m.Outputs)
		writeCodingKeys(w, m.Outputs)
		out(w, "}")
		out(w, "")
		writeDecoderInit(w, m.CamelName+"Output", m.Outputs)

	}

}

// writeFields to writer.
func writeFields(w io.Writer, fields []core.Field) {
	for _, f := range fields {
		out(w, "    // %s", f.Description)
		out(w, "    var %s: %s = %s\n", f.LowerCamelName, swiftType(f), swiftDefault(f))
	}
}

// writeCodingKeys to writer
func writeCodingKeys(w io.Writer, fields []core.Field) {
	if len(fields) == 0 {
		return
	}
	out(w, "    enum CodingKeys: String, CodingKey {")
	for _, f := range fields {
		out(w, "        case %s = \"%s\"", f.LowerCamelName, f.Name)
	}
	out(w, "    }")
}

// writeDecoderInit to writer.
func writeDecoderInit(w io.Writer, extensionName string, fields []core.Field) {
	if len(fields) == 0 {
		return
	}
	out(w, "extension %s {", extensionName)
	out(w, "    init(from decoder: Decoder) throws {")
	out(w, "        let container = try decoder.container(keyedBy: CodingKeys.self)")
	for _, f := range fields {
		out(w, "        if let %s = try container.decodeIfPresent(%s.self, forKey: .%s) {", f.LowerCamelName, swiftType(f), f.LowerCamelName)
		out(w, "            self.%s = %s", f.LowerCamelName, f.LowerCamelName)
		out(w, "        }")
		out(w, "")
	}
	out(w, "    }")
	out(w, "}")
}

// swiftType returns a Swift equivalent type for field f.
func swiftType(f core.Field) string {
	t := f.Type.SwiftType
	if t == "" {
		t = f.Type.CamelName
	}
	if f.IsArray {
		return "[" + t + "]"
	}
	return t
}

func swiftDefault(f core.Field) string {
	s := f.Type.SwiftDefault
	if s == "" {
		s = f.Type.CamelName + "()"
	}
	if f.IsArray {
		return "[]"
	}
	return s
}
