package kotlin

import (
	"io"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateTypesFileConfig struct {
	Output  string
	Types   []core.Type
	Package string
}

func GenerateTypesFile(c GenerateTypesFileConfig) error {
	return common.GenerateFile(c.Output, func(w io.Writer) error {
		common.GenerateWarning(w)
		out(w, "package %s", c.Package)
		out(w, "")
		out(w, "import kotlinx.serialization.SerialName")
		out(w, "import kotlinx.serialization.Serializable")
		out(w, "")
		GenerateTypes(w, c.Types)
		return nil
	})
}

func GenerateTypes(w io.Writer, ts []core.Type) {
	for _, t := range ts {
		out(w, "// %s", t.Description)
		out(w, "@Serializable")
		prefix := ""
		if len(t.Fields) > 0 {
			prefix = "data "
		}
		out(w, "%sclass %s(", prefix, t.CamelName)
		writeFields(w, t.Fields)
		out(w, ")")
		out(w, "")
	}
}

func GenerateMethodTypes(w io.Writer, ms []core.Method) {
	for _, m := range ms {
		prefix := ""
		if len(m.Inputs) > 0 {
			prefix = "data "
		}
		out(w, "@Serializable")
		out(w, "%sclass %sInput(", prefix, m.CamelName)
		writeFields(w, m.Inputs)
		out(w, ")")
		out(w, "")

		prefix = ""
		if len(m.Outputs) > 0 {
			prefix = "data "
		}
		out(w, "@Serializable")
		out(w, "%sclass %sOutput(", prefix, m.CamelName)
		writeFields(w, m.Outputs)
		out(w, ")")
		out(w, "")
	}
}

func writeFields(w io.Writer, fs []core.Field) {
	for _, f := range fs {
		out(w, "    @SerialName(\"%s\") val %s: %s = %s,", f.Name, f.LowerCamelName, kotlinType(f), kotlinDefault(f))
	}
}

func kotlinType(f core.Field) string {
	t := f.Type.KotlinType
	if t == "" {
		t = f.Type.CamelName
	}
	if f.IsArray {
		return "List<" + t + ">"
	}
	return t
}

func kotlinDefault(f core.Field) string {
	s := f.Type.KotlinDefault
	if s == "" {
		s = f.Type.CamelName + "()"
	}
	if f.IsArray {
		return "emptyList()"
	}
	return s
}
