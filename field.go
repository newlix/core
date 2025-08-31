package core

// Field model.
type Field struct {
	Name           string `json:"name"`
	CamelName      string `json:"camel_name"`
	LowerCamelName string `json:"lower_camel_name"`
	Description    string `json:"description"`
	Type           Type   `json:"type"`
	IsArray        bool   `json:"is_array"`
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
