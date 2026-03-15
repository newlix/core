package core

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// Field model.
type Field struct {
	Name           string `json:"name"`
	CamelName      string `json:"camel_name"`
	LowerCamelName string `json:"lower_camel_name"`
	Description    string `json:"description"`
	Type           Type   `json:"type"`
	IsArray        bool   `json:"is_array"`
}

func initFields(ff []Field) ([]Field, error) {
	for i, f := range ff {
		if f.Name == "" {
			return nil, fmt.Errorf("field name must not be empty")
		}
		if f.LowerCamelName == "" {
			f.LowerCamelName = strcase.ToLowerCamel(f.Name)
		}
		if f.CamelName == "" {
			f.CamelName = strcase.ToCamel(f.Name)
		}
		var err error
		f.Type, err = initType(f.Type)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", f.Name, err)
		}
		ff[i] = f
	}
	return ff, nil
}

func BuiltinTypeFields(ff []Field) []Field {
	out := []Field{}
	for _, f := range ff {
		if f.Type.IsBuiltin() {
			out = append(out, f)
		}
	}
	return out
}
