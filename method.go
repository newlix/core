package core

import (
	"log"
	"sort"

	"github.com/iancoleman/strcase"
)

// Method model.
type Method struct {
	Name           string  `json:"name"` // based on filename
	CamelName      string  `json:"camel_name"`
	LowerCamelName string  `json:"lower_camel_name"`
	Description    string  `json:"description"`
	Inputs         []Field `json:"inputs"`
	Outputs        []Field `json:"outputs"`
}

// InitMethods initializes and validates methods.
func InitMethods(mm ...Method) []Method {
	sort.Slice(mm, func(i, j int) bool {
		return mm[i].Name < mm[j].Name
	})
	names := map[string]any{}
	for _, m := range mm {
		if _, ok := names[m.Name]; ok {
			log.Fatalf("duplicate method name = %q", m.Name)
		}
		names[m.Name] = nil
	}

	for i, m := range mm {
		if m.LowerCamelName == "" {
			m.LowerCamelName = strcase.ToLowerCamel(m.Name)
		}
		if m.CamelName == "" {
			m.CamelName = strcase.ToCamel(m.Name)
		}

		m.Inputs = initFields(m.Inputs)
		m.Outputs = initFields(m.Outputs)
		mm[i] = m
	}

	return mm
}
