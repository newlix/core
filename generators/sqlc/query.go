package sqlc

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/newlix/core"
)

type GenerateQueriesFileConfig struct {
	Output  string
	Methods []core.Method
	Types   []core.Type
}

func GenerateQueriesFile(c GenerateQueriesFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	GenerateQueries(w, c.Methods, c.Types)
}

func GenerateQueries(w io.Writer, mm []core.Method, tt []core.Type) {
	generateWarning(w)
	for i, m := range mm {
		writeQuery(w, m, tt)
		if i < len(mm)-1 {
			out(w, "")
		}
	}
}

func writeQuery(w io.Writer, m core.Method, tt []core.Type) {
	outputs := builtinOutputFields(m, tt)
	annotation := queryAnnotation(outputs)
	out(w, "-- name: %s %s", m.CamelName, annotation)

	switch {
	case len(outputs) > 0:
		writeSelect(w, m, outputs, tt)
	case len(m.Inputs) > 0:
		writeInsert(w, m, tt)
	default:
		// no inputs, no outputs — nothing meaningful to generate
	}
}

func writeSelect(w io.Writer, m core.Method, outputs []core.Field, tt []core.Type) {
	table := findTable(outputs[0], tt)
	cols := columnNames(outputs)
	out(w, "SELECT %s FROM %s;", strings.Join(cols, ", "), table)
}

func writeInsert(w io.Writer, m core.Method, tt []core.Type) {
	inputs := expandInputFields(m.Inputs, tt)
	table := findTableForInputs(m.Inputs, tt)
	var cols []string
	var params []string
	for _, f := range inputs {
		cols = append(cols, f.Name)
		params = append(params, sqlcParam(f.Name))
	}
	out(w, "INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(cols, ", "), strings.Join(params, ", "))
}

func queryAnnotation(outputs []core.Field) string {
	switch {
	case len(outputs) == 0:
		return ":exec"
	case hasArray(outputs):
		return ":many"
	default:
		return ":one"
	}
}

func hasArray(ff []core.Field) bool {
	for _, f := range ff {
		if f.IsArray {
			return true
		}
	}
	return false
}

// builtinOutputFields expands composite output fields into their builtin sub-fields.
// If the original field is an array, the expanded sub-fields inherit IsArray.
func builtinOutputFields(m core.Method, tt []core.Type) []core.Field {
	var out []core.Field
	for _, f := range m.Outputs {
		if f.Type.IsBuiltin() {
			out = append(out, f)
		} else {
			t := findType(f.Type.Name, tt)
			for _, sf := range core.BuiltinTypeFields(t.Fields) {
				if f.IsArray {
					sf.IsArray = true
				}
				out = append(out, sf)
			}
		}
	}
	return out
}

// expandInputFields expands composite input fields into their builtin sub-fields.
func expandInputFields(inputs []core.Field, tt []core.Type) []core.Field {
	var out []core.Field
	for _, f := range inputs {
		if f.Type.IsBuiltin() {
			out = append(out, f)
		} else {
			t := findType(f.Type.Name, tt)
			for _, sf := range core.BuiltinTypeFields(t.Fields) {
				out = append(out, sf)
			}
		}
	}
	return out
}

func findType(name string, tt []core.Type) core.Type {
	for _, t := range tt {
		if t.Name == name {
			return t
		}
	}
	log.Fatalf("type %q not found", name)
	return core.Type{}
}

func findTable(f core.Field, tt []core.Type) string {
	if f.Type.IsBuiltin() {
		// look for a type that contains this field
		for _, t := range tt {
			for _, tf := range t.Fields {
				if tf.Name == f.Name {
					return t.Name
				}
			}
		}
	}
	return f.Type.Name
}

func findTableForInputs(inputs []core.Field, tt []core.Type) string {
	for _, f := range inputs {
		if !f.Type.IsBuiltin() {
			return f.Type.Name
		}
	}
	// fallback: find a type that has matching field names
	for _, f := range inputs {
		for _, t := range tt {
			for _, tf := range t.Fields {
				if tf.Name == f.Name {
					return t.Name
				}
			}
		}
	}
	log.Fatal("cannot determine table for inputs")
	return ""
}

func columnNames(ff []core.Field) []string {
	var out []string
	for _, f := range ff {
		out = append(out, f.Name)
	}
	return out
}

func sqlcParam(name string) string {
	return "@" + name
}
