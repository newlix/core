package kotlin

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateTypesFileConfig struct {
	Output  string
	Types   []core.Type
	Package string
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

	out(w, "package %s\n\n", c.Package)

	out(w, `
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
`)
	GenerateTypes(w, c.Types)
}

func GenerateTypes(w io.Writer, ts []core.Type) {
	// types
	for _, t := range ts {
		out(w, "// %s", t.Description)
		out(w, "@Serializable")
		if len(t.Fields) > 0 {
			fmt.Fprint(w, "data ")
		}
		out(w, "class %s(", t.CamelName)
		writeFields(w, t.Fields)
		out(w, ")")
		out(w, "")
	}

}

func GenerateMethodTypes(w io.Writer, ms []core.Method) {
	for _, m := range ms {
		// inputs
		out(w, "@Serializable")
		if len(m.Inputs) > 0 {
			fmt.Fprint(w, "data ")
		}
		out(w, "class %sInput(", m.CamelName)
		writeFields(w, m.Inputs)
		out(w, ")")
		out(w, "")

		// outputs
		out(w, "@Serializable")
		if len(m.Outputs) > 0 {
			fmt.Fprint(w, "data ")
		}
		out(w, "class %sOutput(", m.CamelName)
		writeFields(w, m.Outputs)
		out(w, ")")
		out(w, "")
	}

}

// writeFields to writer.
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
