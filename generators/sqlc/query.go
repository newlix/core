package sqlc

import (
	"fmt"
	"io"
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

func GenerateQueriesFile(c GenerateQueriesFileConfig) error {
	if err := os.MkdirAll(path.Dir(c.Output), 0o700); err != nil {
		return err
	}
	w, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer w.Close()

	return GenerateQueries(w, c.Methods, c.Types)
}

func GenerateQueries(w io.Writer, mm []core.Method, tt []core.Type) error {
	generateWarning(w)
	for i, m := range mm {
		if err := writeQuery(w, m, tt); err != nil {
			return err
		}
		if i < len(mm)-1 {
			out(w, "")
		}
	}
	return nil
}

func writeQuery(w io.Writer, m core.Method, tt []core.Type) error {
	outputs, err := builtinOutputFields(m, tt)
	if err != nil {
		return err
	}
	annotation := queryAnnotation(outputs)
	out(w, "-- name: %s %s", m.CamelName, annotation)

	switch {
	case len(outputs) > 0:
		return writeSelect(w, m, outputs, tt)
	case len(m.Inputs) > 0:
		return writeInsert(w, m, tt)
	default:
		// no inputs, no outputs — nothing meaningful to generate
		return nil
	}
}

func writeSelect(w io.Writer, m core.Method, outputs []core.Field, tt []core.Type) error {
	table := m.Table
	if table == "" {
		var err error
		table, err = findTable(outputs[0], tt)
		if err != nil {
			return err
		}
	}
	cols := columnNames(outputs)
	out(w, "SELECT %s FROM %s;", strings.Join(cols, ", "), table)
	return nil
}

func writeInsert(w io.Writer, m core.Method, tt []core.Type) error {
	inputs, err := expandInputFields(m.Inputs, tt)
	if err != nil {
		return err
	}
	table := m.Table
	if table == "" {
		table, err = findTableForInputs(m.Inputs, tt)
		if err != nil {
			return err
		}
	}
	var cols []string
	var params []string
	for _, f := range inputs {
		cols = append(cols, f.Name)
		params = append(params, sqlcParam(f.Name))
	}
	out(w, "INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(cols, ", "), strings.Join(params, ", "))
	return nil
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
func builtinOutputFields(m core.Method, tt []core.Type) ([]core.Field, error) {
	var fields []core.Field
	for _, f := range m.Outputs {
		if f.Type.IsBuiltin() {
			fields = append(fields, f)
		} else {
			t, err := findType(f.Type.Name, tt)
			if err != nil {
				return nil, err
			}
			for _, sf := range core.BuiltinTypeFields(t.Fields) {
				if f.IsArray {
					sf.IsArray = true
				}
				fields = append(fields, sf)
			}
		}
	}
	return fields, nil
}

// expandInputFields expands composite input fields into their builtin sub-fields.
func expandInputFields(inputs []core.Field, tt []core.Type) ([]core.Field, error) {
	var fields []core.Field
	for _, f := range inputs {
		if f.Type.IsBuiltin() {
			fields = append(fields, f)
		} else {
			t, err := findType(f.Type.Name, tt)
			if err != nil {
				return nil, err
			}
			for _, sf := range core.BuiltinTypeFields(t.Fields) {
				fields = append(fields, sf)
			}
		}
	}
	return fields, nil
}

func findType(name string, tt []core.Type) (core.Type, error) {
	for _, t := range tt {
		if t.Name == name {
			return t, nil
		}
	}
	return core.Type{}, fmt.Errorf("type %q not found", name)
}

func findTable(f core.Field, tt []core.Type) (string, error) {
	if f.Type.IsBuiltin() {
		// look for a type that contains this field
		for _, t := range tt {
			for _, tf := range t.Fields {
				if tf.Name == f.Name {
					return t.Name, nil
				}
			}
		}
		return "", fmt.Errorf("cannot determine table for field %q", f.Name)
	}
	return f.Type.Name, nil
}

func findTableForInputs(inputs []core.Field, tt []core.Type) (string, error) {
	for _, f := range inputs {
		if !f.Type.IsBuiltin() {
			return f.Type.Name, nil
		}
	}
	// fallback: find a type that has matching field names
	for _, f := range inputs {
		for _, t := range tt {
			for _, tf := range t.Fields {
				if tf.Name == f.Name {
					return t.Name, nil
				}
			}
		}
	}
	return "", fmt.Errorf("cannot determine table for inputs")
}

func columnNames(ff []core.Field) []string {
	var cols []string
	for _, f := range ff {
		cols = append(cols, f.Name)
	}
	return cols
}

func sqlcParam(name string) string {
	return "@" + name
}
