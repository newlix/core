package golang

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateClientFileConfig struct {
	Output  string
	Package string
	Methods []core.Method
	Types   []core.Type
}

// generate implementation.
func GenerateClientFile(c GenerateClientFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	common.GenerateWarning(w)
	out(w, "package %s", PackageName(c.Package))
	out(w, `import (
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
)

`)
	GenerateImports(w, c.Package, c.Types)
	GenerateMethodTypes(w, c.Package, c.Methods, c.Types)

	GenerateClient(w, c.Methods)

}

func GenerateClient(w io.Writer, mm []core.Method) {
	w.Write([]byte(`// Client is the API client.
type Client struct {
	// URL is the required API endpoint address.
	URL string
	// AuthToken is an optional authentication token.
	AuthToken string
	// HTTPClient is the client used for making requests, defaulting to http.DefaultClient.
	HTTPClient *http.Client
}

// Error is an error returned by the client.
type Error struct {
	Status     string
	StatusCode int
	Type       string
	Message    string
}

// Error implementation.
func (e Error) Error() string {
	if e.Type == "" {
		return fmt.Sprintf("%s: %d", e.Status, e.StatusCode)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// call implementation.
func call(client *http.Client, authToken, endpoint, method string, in, out interface{}) error {
	var body io.Reader

	// default client
	if client == nil {
		client = http.DefaultClient
	}

	// input params
	if in != nil {
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(in)
		if err != nil {
			return fmt.Errorf("encoding: %w", err)
		}
		body = &buf
	}

	// POST request
	req, err := http.NewRequest("POST", endpoint+"/"+method, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// auth token
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	// response
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// error
	if res.StatusCode >= 300 {
		var e Error
		if res.Header.Get("Content-Type") == "application/json" {
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				return err
			}
		}
		e.Status = http.StatusText(res.StatusCode)
		e.StatusCode = res.StatusCode
		return e
	}

	// output params
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return err
		}
	}

	return nil
}

`))

	for i, m := range mm {
		out(w, "// %s", m.Description)
		out(w, "func (c *Client) %s(in %sInput) (%sOutput, error) {", m.CamelName, m.CamelName, m.CamelName)
		out(w, "	var out %sOutput", m.CamelName)
		out(w, "	return out, call(c.HTTPClient, c.AuthToken, c.URL, %q, in, &out)", m.Name)
		out(w, "}")
		if i < len(mm)-1 {
			out(w, "")
		}
	}
}
