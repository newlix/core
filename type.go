package core

import (
	"log"
	"sort"

	"github.com/iancoleman/strcase"
)

// Type model.
type Type struct {
	Name           string   `json:"name"` // based on filename
	CamelName      string   `json:"camel_name"`
	LowerCamelName string   `json:"lower_camel_name"`
	Description    string   `json:"description"`
	Fields         []Field  `json:"fields"`
	PrimaryKey     []string `json:"primary_key"`
	GoPackage      string   `json:"go_import"`
	GoType         string   `json:"go_type"`
	SwiftType      string   `json:"swift_type"`
	SwiftDefault   string   `json:"swift_default"`
	KotlinType     string   `json:"kotlin_type"`
	KotlinDefault  string   `json:"kotlin_default"`
	SqlcType       string   `json:"sqlc_type"`
	isInitialized  bool
	isBuiltin      bool
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
	names := map[string]struct{}{}
	for _, t := range tt {
		if _, ok := names[t.Name]; ok {
			log.Fatalf("duplicate type name = %q", t.Name)
		}
		names[t.Name] = struct{}{}
	}

	for i, t := range tt {
		tt[i] = initType(t)
	}
	return tt
}

func initType(t Type) Type {
	if t.LowerCamelName == "" {
		t.LowerCamelName = strcase.ToLowerCamel(t.Name)
	}
	if t.CamelName == "" {
		t.CamelName = strcase.ToCamel(t.Name)
	}
	if len(t.PrimaryKey) == 0 {
		t.PrimaryKey = []string{"id"}
	}
	t.Fields = initFields(t.Fields)
	t.isInitialized = true
	return t
}

var String = Type{
	Name:          "string",
	Description:   "built-in string",
	GoType:        "string",
	SwiftType:     "String",
	SwiftDefault:  `""`,
	KotlinType:    "String",
	KotlinDefault: `""`,
	SqlcType:      "TEXT",
	isBuiltin:     true,
}

var Int = Type{
	Name:          "int",
	Description:   "built-in integer",
	GoType:        "int",
	SwiftType:     "Int",
	SwiftDefault:  "0",
	KotlinType:    "Long",
	KotlinDefault: "0",
	SqlcType:      "BIGINT",
	isBuiltin:     true,
}

var Bool = Type{
	Name:          "bool",
	Description:   "built-in bool",
	GoType:        "bool",
	SwiftType:     "Bool",
	SwiftDefault:  "false",
	KotlinType:    "Boolean",
	KotlinDefault: "false",
	SqlcType:      "BOOLEAN",
	isBuiltin:     true,
}

var Float = Type{
	Name:          "float",
	Description:   "built-in float",
	GoType:        "float64",
	SwiftType:     "Double",
	SwiftDefault:  "0.0",
	KotlinType:    "Double",
	KotlinDefault: "0.0",
	SqlcType:      "DOUBLE PRECISION",
	isBuiltin:     true,
}

var Time = Type{
	Name:          "time",
	Description:   "built-in time",
	GoType:        "*time.Time",
	SwiftType:     "Date",
	SwiftDefault:  "Date()",
	KotlinType:    "Date",
	KotlinDefault: "Date()",
	SqlcType:      "TIMESTAMPTZ",
	isBuiltin:     true,
}
