package golang

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

// goInitialisms maps uppercase initialisms recognized by Go conventions.
var goInitialisms = map[string]bool{
	"ID": true, "URL": true, "HTTP": true, "API": true,
	"JSON": true, "XML": true, "SQL": true, "HTML": true,
	"CSS": true, "IP": true, "TCP": true, "UDP": true,
	"RPC": true, "SSH": true, "TLS": true, "TTL": true,
	"UUID": true, "ASCII": true, "UTF8": true, "ACL": true,
	"EOF": true, "QPS": true, "DNS": true,
}

var camelWordRe = regexp.MustCompile(`[A-Z][a-z]*`)

// GoName converts a CamelCase name to Go-idiomatic form by uppercasing initialisms.
// e.g. "UserId" → "UserID", "HttpUrl" → "HTTPURL", "ApiKey" → "APIKey"
func GoName(name string) string {
	return camelWordRe.ReplaceAllStringFunc(name, func(word string) string {
		upper := strings.ToUpper(word)
		if goInitialisms[upper] {
			return upper
		}
		return word
	})
}

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

	if err := os.MkdirAll(path.Dir(c.Output), 0o700); err != nil {
		log.Fatal(err)
	}
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)

	out(w, "package %s", PackageName(c.Package))
	out(w, `import "encoding/json"`)
	needsTime := false
	for _, t := range c.Types {
		for _, f := range t.Fields {
			if f.Type.Name == core.Time.Name {
				needsTime = true
				break
			}
		}
		if needsTime {
			break
		}
	}
	if needsTime {
		out(w, `import "time"`)
	}
	out(w, `
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
	imports := make([]string, 0, len(m))
	for s := range m {
		imports = append(imports, s)
	}
	sort.Strings(imports)
	for _, s := range imports {
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
		out(w, "type %s struct {", GoName(t.CamelName))
		writeFields(w, pkg, t.Fields, tags)
		out(w, "}")

		GenerateJSONMarshalerIfNeeded(w, t, pkg, tags)

	}
}

func GenerateJSONMarshalerIfNeeded(w io.Writer, t core.Type, pkg string, tags []string) {
	needsCustomJSON := false
	for _, f := range t.Fields {
		if f.Type.GoType == "*time.Time" {
			needsCustomJSON = true
			break
		}
	}
	if !needsCustomJSON {
		return
	}

	name := GoName(t.CamelName)
	out(w, "")
	out(w, "func (p %s) MarshalJSON() ([]byte, error) {", name)
	writeAliasType(w, pkg, t, tags)
	out(w, "	return json.Marshal(&Alias{")
	for _, f := range t.Fields {
		fn := GoName(f.CamelName)
		if f.Type.GoType == "*time.Time" {
			out(w, "		%s: func() int64 { if p.%s != nil { return p.%s.Unix() }; return 0 }(),", fn, fn, fn)
		} else {
			out(w, "		%s: p.%s,", fn, fn)
		}
	}
	out(w, "	})")
	out(w, "}")
	out(w, "")

	out(w, "func (p *%s) UnmarshalJSON(data []byte) error {", name)
	writeAliasType(w, pkg, t, tags)
	out(w, "	var tmp Alias")
	out(w, "	if err := json.Unmarshal(data, &tmp); err != nil {")
	out(w, "		return err")
	out(w, "	}")
	for _, f := range t.Fields {
		fn := GoName(f.CamelName)
		if f.Type.GoType == "*time.Time" {
			out(w, "	if tmp.%s != 0 {", fn)
			out(w, "		p.%s = ptr(time.Unix(tmp.%s, 0))", fn, fn)
			out(w, "	}")
		} else {
			out(w, "	p.%s = tmp.%s", fn, fn)
		}
	}
	out(w, "	return nil")
	out(w, "}")
}
func GenerateMethodTypes(w io.Writer, pkg string, mm []core.Method) {
	tags := []string{"json"}

	// methods
	for i, m := range mm {
		name := GoName(m.CamelName)
		out(w, "type %sInput struct {", name)
		writeFields(w, pkg, m.Inputs, tags)
		out(w, "}")
		out(w, "")
		out(w, "type %sOutput struct {", name)
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
		out(w, "	%s %s %s", GoName(f.CamelName), FieldGoType(pkg, f, false), goTags(f.Name, tags))
		if i < len(ff)-1 {
			fmt.Fprintf(w, "\n")
		}
	}
}

func writeAliasFields(w io.Writer, pkg string, ff []core.Field, tags []string) {
	for _, f := range ff {
		out(w, "		%s %s %s", GoName(f.CamelName), FieldGoType(pkg, f, true), goTags(f.Name, tags))
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
