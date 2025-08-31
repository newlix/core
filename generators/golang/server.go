package golang

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateServerFileConfig struct {
	Output  string
	Package string
	Methods []core.Method
	Types   []core.Type
}

// generate implementation.
func GenerateServerFile(c GenerateServerFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)

	out(w, "package %s", PackageName(c.Package))
	out(w, "")
	out(w, `import (
	"net/http"
	"github.com/newlix/core"
)`)
	GenerateImports(w, c.Package, c.Types)
	GenerateMethodTypes(w, c.Package, c.Methods, c.Types)

	GenerateServer(w, c.Methods)

}

// Generate writes the Go server implementations to w.
func GenerateServer(w io.Writer, mm []core.Method) {
	out(w, "// ServeHTTP implementation.")
	out(w, "func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {")
	out(w, "	if r.Method != \"POST\" {")
	out(w, "		w.Header().Add(\"Allow\", \"POST\") // RFC 9110.")
	out(w, "		core.WriteError(w, core.BadRequest(\"only POST is allowed\"))")
	out(w, "		return")
	out(w, "	}")
	out(w, "	ctx := core.NewRequestContext(r.Context(), r)")
	out(w, "	var res interface{}")
	out(w, "	var err error")
	out(w, "	switch r.URL.Path {")

	for _, m := range mm {

		out(w, "	case \"/%s\":", m.Name)
		// parse input
		out(w, "		var in %sInput", m.CamelName)
		out(w, "		var out %sOutput", m.CamelName)
		out(w, "		err = core.ReadRequest(r, &in)")
		out(w, "		if err != nil {")
		out(w, "			break")
		out(w, "		}")
		out(w, "		out, err = s.%s(ctx, in)", m.CamelName)
		out(w, "		res = out")
	}
	out(w, "	default:")
	out(w, "		err = core.BadRequest(\"Invalid method\")")
	out(w, "	}")
	out(w, "")
	out(w, "	if err != nil {")
	out(w, "		core.WriteError(w, err)")
	out(w, "		return")
	out(w, "	}")
	out(w, "")
	out(w, "	core.WriteResponse(w, res)")
	out(w, "	return")
	out(w, "}")
}
