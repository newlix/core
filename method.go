package core

import (
	"fmt"
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
func InitMethods(mm ...Method) ([]Method, error) {
	sort.Slice(mm, func(i, j int) bool {
		return mm[i].Name < mm[j].Name
	})
	names := map[string]struct{}{}
	for _, m := range mm {
		if m.Name == "" {
			return nil, fmt.Errorf("method name must not be empty")
		}
		if _, ok := names[m.Name]; ok {
			return nil, fmt.Errorf("duplicate method name = %q", m.Name)
		}
		names[m.Name] = struct{}{}
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

	return mm, nil
}
