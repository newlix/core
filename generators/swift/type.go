package swift

import (
	"io"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateTypesFileConfig struct {
	Output string
	Types  []core.Type
}

func GenerateTypesFile(c GenerateTypesFileConfig) error {
	if err := os.MkdirAll(path.Dir(c.Output), 0o700); err != nil {
		return err
	}
	w, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer w.Close()

	common.GenerateWarning(w)
	out(w, "import Foundation")
	out(w, "")

	GenerateTypes(w, c.Types)
	return nil
}

func GenerateTypes(w io.Writer, tt []core.Type) {
	for i, t := range tt {
		out(w, "// %s", t.Description)
		out(w, "struct %s: Codable {", t.CamelName)
		writeFields(w, t.Fields)
		writeCodingKeys(w, t.Fields)
		out(w, "}")
		writeDecoderInit(w, t.CamelName, t.Fields)
		if i < len(tt)-1 {
			out(w, "")
		}
	}
}

func GenerateMethodTypes(w io.Writer, ms []core.Method) {
	for i, m := range ms {
		out(w, "struct %sInput: Codable {", m.CamelName)
		writeFields(w, m.Inputs)
		writeCodingKeys(w, m.Inputs)
		out(w, "}")
		writeDecoderInit(w, m.CamelName+"Input", m.Inputs)

		out(w, "")

		out(w, "struct %sOutput: Codable {", m.CamelName)
		writeFields(w, m.Outputs)
		writeCodingKeys(w, m.Outputs)
		out(w, "}")
		writeDecoderInit(w, m.CamelName+"Output", m.Outputs)

		if i < len(ms)-1 {
			out(w, "")
		}
	}
}

// writeFields to writer.
func writeFields(w io.Writer, fields []core.Field) {
	for i, f := range fields {
		out(w, "    // %s", f.Description)
		out(w, "    var %s: %s = %s", f.LowerCamelName, swiftType(f), swiftDefault(f))
		if i < len(fields)-1 {
			out(w, "")
		}
	}
}

// writeCodingKeys to writer.
func writeCodingKeys(w io.Writer, fields []core.Field) {
	if len(fields) == 0 {
		return
	}
	out(w, "")
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
	out(w, "")
	out(w, "extension %s {", extensionName)
	out(w, "    init(from decoder: Decoder) throws {")
	out(w, "        let container = try decoder.container(keyedBy: CodingKeys.self)")
	for i, f := range fields {
		out(w, "        if let %s = try container.decodeIfPresent(%s.self, forKey: .%s) {", f.LowerCamelName, swiftType(f), f.LowerCamelName)
		out(w, "            self.%s = %s", f.LowerCamelName, f.LowerCamelName)
		out(w, "        }")
		if i < len(fields)-1 {
			out(w, "")
		}
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
