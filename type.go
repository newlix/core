package core

import (
	"log"
	"sort"

	"github.com/iancoleman/strcase"
)

// Type model.
type Type struct {
	Name            string   `json:"name"` // based on filename
	CamelName       string   `json:"camel_name"`
	LowerCamelName  string   `json:"lower_camel_name"`
	Description     string   `json:"description"`
	Fields          []Field  `json:"fields"`
	PrimaryKey      []string `json:"primary_key"`
	GoPackage       string   `json:"go_import"`
	GoType          string   `json:"go_type"`
	SwiftType       string   `json:"swift_type"`
	SwiftDefault    string   `json:"swift_default"`
	KotlinType      string   `json:"kotlin_type"`
	KotlinDefault string `json:"kotlin_default"`
	isInitialized bool
	isBuiltin       bool
}

func (t Type) IsInitialized() bool {
	return t.isInitialized
}

func (t Type) IsBuiltin() bool {
	return t.isBuiltin
}

func InitTypes(tt ...Type) []Type {
	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Name < tt[j].Name
	})
	names := map[string]any{}
	for _, t := range tt {
		if _, ok := names[t.Name]; ok {
			log.Fatalf("duplicate type name = %q", t.Name)
		}
		names[t.Name] = nil
	}

	for i, t := range tt {
		// overwrite
		if t.LowerCamelName == "" {
			t.LowerCamelName = strcase.ToLowerCamel(t.Name)
		}

		if t.CamelName == "" {
			t.CamelName = strcase.ToCamel(t.Name)
		}

		if len(t.PrimaryKey) == 0 {
			t.PrimaryKey = []string{"id"}
		}

		for j, f := range t.Fields {

			if f.LowerCamelName == "" {
				f.LowerCamelName = strcase.ToLowerCamel(f.Name)
			}
			if f.CamelName == "" {
				f.CamelName = strcase.ToCamel(f.Name)
			}
			f.Type = InitTypes(f.Type)[0]
			t.Fields[j] = f
		}
		t.isInitialized = true
		tt[i] = t
	}
	return tt
}

var String = Type{
	Name:            "string",
	Description:     "built-in string",
	Fields:          nil,
	PrimaryKey:      nil,
	GoType:          "string",
	SwiftType:       "String",
	SwiftDefault:  `""`,
	KotlinType:    "String",
	KotlinDefault: `""`,
	isBuiltin:     true,
}

var Int = Type{
	Name:            "int",
	Description:     "built-in integer",
	Fields:          nil,
	PrimaryKey:      nil,
	GoType:          "int",
	SwiftType:       "Int",
	SwiftDefault:    "0",
	KotlinType:      "Long",
	KotlinDefault: "0",
	isBuiltin:     true,
}

var Bool = Type{
	Name:            "bool",
	Description:     "built-in bool",
	Fields:          nil,
	PrimaryKey:      nil,
	GoType:          "bool",
	SwiftType:       "Bool",
	SwiftDefault:    "false",
	KotlinType:      "Boolean",
	KotlinDefault: "false",
	isBuiltin:     true,
}

var Float = Type{
	Name:            "float",
	Description:     "built-in float",
	Fields:          nil,
	PrimaryKey:      nil,
	GoType:          "float64",
	SwiftType:       "Double",
	SwiftDefault:    "0.0",
	KotlinType:      "Double",
	KotlinDefault: "0.0",
	isBuiltin:     true,
}

var Time = Type{
	Name:            "time",
	Description:     "built-in date",
	Fields:          nil,
	PrimaryKey:      nil,
	GoType:          "*time.Time",
	SwiftType:       "Date",
	SwiftDefault:    "Date()",
	KotlinType:      "Date",
	KotlinDefault: "Date()",
	isBuiltin:     true,
}
