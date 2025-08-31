package golang

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

var DefaultTags = []string{"json", "db"}

type GenerateTypesFileConfig struct {
	Output  string
	Package string
	Types   []core.Type
	Tags    []string
}

func GenerateTypesFile(c GenerateTypesFileConfig) {
	if len(c.Tags) == 0 {
		c.Tags = DefaultTags
	}
	for _, t := range c.Types {
		if !t.IsInitialized() {
			log.Fatalf("type %s is not initialized", t.Name)
		}
	}

	CheckPackage(c.Package, c.Types)

	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)

	out(w, "package %s", PackageName(c.Package))
	out(w, `import (
		"encoding/json"
		"time"
)

func ptr[T any](x T) *T {
    return &x
}
`)
	GenerateImports(w, c.Package, c.Types)

	GenerateTypes(w, c.Package, c.Types, c.Tags)

}

func CheckPackage(pkg string, tt []core.Type) {
	for _, t := range tt {
		if t.GoPackage != pkg {
			log.Fatalf("%s belongs to %q not %q", t.Name, t.GoPackage, pkg)
		}
	}
}

func GenerateImports(w io.Writer, pkg string, tt []core.Type) {
	m := map[string]bool{}
	for _, t := range tt {
		if t.GoPackage != pkg {
			m[t.GoPackage] = true
		}
		for _, f := range t.Fields {
			if f.Type.GoPackage != "" {
				if f.Type.GoPackage != pkg {
					m[f.Type.GoPackage] = true
				}
			}
		}
	}
	for s := range m {
		if strings.HasSuffix(s, `"`) {
			out(w, "import %s", s)
		} else {
			out(w, "import %q", s)
		}
	}
}

func GenerateTypes(w io.Writer, pkg string, tt []core.Type, tags []string) {
	for _, t := range tt {
		out(w, "// %s", t.Description)
		out(w, "type %s struct {", t.CamelName)
		writeFields(w, pkg, t.Fields, tags)
		out(w, "}")

		GenerateJSONMashlerIfNeeded(w, t, pkg, tags)

	}
}

func GenerateJSONMashlerIfNeeded(w io.Writer, t core.Type, pkg string, tags []string) {
	out(w, "")
	out(w, "func (p %s) MarshalJSON() ([]byte, error) {", t.CamelName)
	writeAliasType(w, pkg, t, tags)
	out(w, "	return json.Marshal(&Alias{")
	for _, f := range t.Fields {
		if f.Type.GoType == "*time.Time" {
			out(w, "		%s: p.%s.Unix(),", f.CamelName, f.CamelName)
		} else {
			out(w, "		%s: p.%s,", f.CamelName, f.CamelName)
		}
	}
	out(w, "	})")
	out(w, "}")
	out(w, "")

	out(w, "func (p *%s) UnmarshalJSON(data []byte) error {", t.CamelName)
	writeAliasType(w, pkg, t, tags)
	out(w, "	var tmp Alias")
	out(w, "	if err := json.Unmarshal(data, &tmp); err != nil {")
	out(w, "		return err")
	out(w, "	}")
	for _, f := range t.Fields {
		if f.Type.GoType == "*time.Time" {
			out(w, "	p.%s = ptr(time.Unix(tmp.%s, 0))", f.CamelName, f.CamelName)
		} else {
			out(w, "	p.%s = tmp.%s", f.CamelName, f.CamelName)
		}
	}
	out(w, "	return nil")
	out(w, "}")
}
func GenerateMethodTypes(w io.Writer, pkg string, mm []core.Method, tt []core.Type) {
	tags := []string{"json"}

	// methods
	for i, m := range mm {
		out(w, "type %sInput struct {", m.CamelName)
		writeFields(w, pkg, m.Inputs, tags)
		out(w, "}")
		out(w, "")
		out(w, "type %sOutput struct {", m.CamelName)
		writeFields(w, pkg, m.Outputs, tags)
		out(w, "}")
		if i < len(mm)-1 {
			out(w, "")
		}
	}
}

func writeAliasType(w io.Writer, pkg string, t core.Type, tags []string) {
	out(w, "	type Alias struct {")
	writeAliasFields(w, pkg, t.Fields, tags)
	out(w, "	}")
}

// writeFields to writer
func writeFields(w io.Writer, pkg string, ff []core.Field, tags []string) {
	for i, f := range ff {
		out(w, "	// %s", f.Description)
		out(w, "	%s %s %s", f.CamelName, FieldGoType(pkg, f, false), goTags(f.Name, tags))
		if i < len(ff)-1 {
			fmt.Fprintf(w, "\n")
		}
	}
}

func writeAliasFields(w io.Writer, pkg string, ff []core.Field, tags []string) {
	for _, f := range ff {
		out(w, "		%s %s %s", f.CamelName, FieldGoType(pkg, f, true), goTags(f.Name, tags))
	}
}

func PackageName(pkg string) string {
	if strings.Contains(pkg, " ") {
		return strings.Split(pkg, " ")[0]
	}
	return path.Base(pkg)
}

func FieldGoType(pkg string, f core.Field, int64Time bool) string {
	if int64Time && (f.Type.GoType == "*time.Time") {
		return "int64"
	}
	if f.IsArray {
		return "[]" + TypeGoType(pkg, f.Type)
	}
	return TypeGoType(pkg, f.Type)
}
func TypeGoType(pkg string, t core.Type) string {
	s := t.GoType
	if s == "" {
		s = t.CamelName
		if t.GoPackage != pkg {
			s = PackageName(t.GoPackage) + "." + s
		}
	}
	return s
}

// goTags returns tags for a field.
func goTags(name string, tags []string) string {
	var s []string

	for _, tag := range tags {
		s = append(s, fmt.Sprintf("%s:%q", tag, name))
	}
	return fmt.Sprintf("`%s`", strings.Join(s, " "))
}
