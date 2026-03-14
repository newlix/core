package core

import "github.com/iancoleman/strcase"

// Field model.
type Field struct {
	Name           string `json:"name"`
	CamelName      string `json:"camel_name"`
	LowerCamelName string `json:"lower_camel_name"`
	Description    string `json:"description"`
	Type           Type   `json:"type"`
	IsArray        bool   `json:"is_array"`
}

func initFields(ff []Field) []Field {
	for i, f := range ff {
		if f.LowerCamelName == "" {
			f.LowerCamelName = strcase.ToLowerCamel(f.Name)
		}
		if f.CamelName == "" {
			f.CamelName = strcase.ToCamel(f.Name)
		}
		f.Type = initType(f.Type)
		ff[i] = f
	}
	return ff
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
