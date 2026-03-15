package core

import (
	"fmt"
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

func InitTypes(tt ...Type) ([]Type, error) {
	sort.Slice(tt, func(i, j int) bool {
		return tt[i].Name < tt[j].Name
	})
	names := map[string]struct{}{}
	for _, t := range tt {
		if t.Name == "" {
			return nil, fmt.Errorf("type name must not be empty")
		}
		if _, ok := names[t.Name]; ok {
			return nil, fmt.Errorf("duplicate type name = %q", t.Name)
		}
		names[t.Name] = struct{}{}
	}

	for i, t := range tt {
		var err error
		tt[i], err = initType(t)
		if err != nil {
			return nil, fmt.Errorf("type %q: %w", t.Name, err)
		}
	}
	return tt, nil
}

func initType(t Type) (Type, error) {
	if t.isInitialized {
		return t, nil
	}
	if t.LowerCamelName == "" {
		t.LowerCamelName = strcase.ToLowerCamel(t.Name)
	}
	if t.CamelName == "" {
		t.CamelName = strcase.ToCamel(t.Name)
	}
	if len(t.PrimaryKey) == 0 {
		t.PrimaryKey = []string{"id"}
	}
	var err error
	t.Fields, err = initFields(t.Fields)
	if err != nil {
		return Type{}, err
	}
	t.isInitialized = true
	return t, nil
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

