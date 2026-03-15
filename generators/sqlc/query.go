// Package sqlc generates sqlc-compatible SQL scaffolding from core types and methods.
//
// The generated queries are starting points — not production-ready SQL.
// Methods with both inputs and outputs generate SELECT without WHERE clauses;
// users must add filtering, joins, and other query logic manually.
package sqlc

import (
	"fmt"
	"io"
	"strings"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateQueriesFileConfig struct {
	Output  string
	Methods []core.Method
	Types   []core.Type
}

func GenerateQueriesFile(c GenerateQueriesFileConfig) error {
	return common.GenerateFile(c.Output, func(w io.Writer) error {
		return GenerateQueries(w, c.Methods, c.Types)
	})
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
	outputs, err := expandFields(m.Outputs, tt, true)
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
	inputs, err := expandFields(m.Inputs, tt, false)
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
	cols := make([]string, 0, len(inputs))
	params := make([]string, 0, len(inputs))
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

// expandFields expands composite fields into their builtin sub-fields.
// If propagateArray is true and the parent field is an array, sub-fields inherit IsArray.
func expandFields(ff []core.Field, tt []core.Type, propagateArray bool) ([]core.Field, error) {
	var fields []core.Field
	for _, f := range ff {
		if f.Type.IsBuiltin() {
			fields = append(fields, f)
			continue
		}
		t, err := findType(f.Type.Name, tt)
		if err != nil {
			return nil, err
		}
		subs := core.BuiltinTypeFields(t.Fields)
		if propagateArray && f.IsArray {
			for i := range subs {
				subs[i].IsArray = true
			}
		}
		fields = append(fields, subs...)
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
	cols := make([]string, 0, len(ff))
	for _, f := range ff {
		cols = append(cols, f.Name)
	}
	return cols
}

func sqlcParam(name string) string {
	return "@" + name
}
