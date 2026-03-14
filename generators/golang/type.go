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
	out(w, "")
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
	if len(imports) > 0 {
		out(w, "")
	}
}

func GenerateTypes(w io.Writer, pkg string, tt []core.Type, tags []string) {
	for i, t := range tt {
		out(w, "// %s", t.Description)
		out(w, "type %s struct {", GoName(t.CamelName))
		writeFields(w, pkg, t.Fields, tags)
		out(w, "}")
		if i < len(tt)-1 {
			out(w, "")
		}
	}
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

// writeFields to writer
func writeFields(w io.Writer, pkg string, ff []core.Field, tags []string) {
	for i, f := range ff {
		out(w, "	// %s", f.Description)
		out(w, "	%s %s %s", GoName(f.CamelName), FieldGoType(pkg, f), goTags(f.Name, tags))
		if i < len(ff)-1 {
			out(w, "")
		}
	}
}

func PackageName(pkg string) string {
	if strings.Contains(pkg, " ") {
		return strings.Split(pkg, " ")[0]
	}
	return path.Base(pkg)
}

func FieldGoType(pkg string, f core.Field) string {
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
