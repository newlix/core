package swift

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
	Methods []core.Method
	Client  string
}

func GenerateClientFile(c GenerateClientFileConfig) {
	if err := os.MkdirAll(path.Dir(c.Output), 0o700); err != nil {
		log.Fatal(err)
	}
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)
	out(w, "import Foundation")
	out(w, "")

	GenerateClient(w, c.Methods, c.Client)
	GenerateMethodTypes(w, c.Methods)

}

const start = `struct CoreError: LocalizedError {
    let status: Int
    let message: String

    var errorDescription: String? {
        message
    }
}

// %s is the API client.
struct %s {
    // encoder is the conventional json encoder
    private let encoder = JSONEncoder()

    // decoder is the conventional json decoder
    private let decoder = JSONDecoder()

    // endpoint is the required API endpoint address.
    let endpoint: String

    // AuthToken is an optional authentication token.
    var authToken: String?

    // session is the client used for making requests, defaulting to URLSession.shared.
    let session: URLSession = URLSession.shared

    private func call<Input, Output>(method: String, input: Input) async throws -> Output where Input: Codable, Output: Codable {
        guard let url = URL(string: endpoint + "/" + method) else {
            throw CoreError(status: 0, message: "Invalid URL: \(endpoint)/\(method)")
        }

        let body = try self.encoder.encode(input)

        var req = URLRequest(url: url)
        req.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let tok = self.authToken {
            req.setValue("Bearer " + tok, forHTTPHeaderField: "Authorization")
        }
        req.httpMethod = "POST"
        req.httpBody = body

        let (data, res) = try await self.session.data(for: req)

        guard let r = res as? HTTPURLResponse else {
            throw CoreError(status: 0, message: "Unexpected response type")
        }

        if r.statusCode >= 300 {
            let body = String(decoding: data, as: UTF8.self)
            let err = CoreError(status: r.statusCode, message: body)
            throw err
        }
        return try self.decoder.decode(Output.self, from: data)
    }
`

// GenerateClient writes the Swift client implementation to w.
func GenerateClient(w io.Writer, ms []core.Method, client string) {
	out(w, start, client, client)
	for _, m := range ms {
		writeMethod(w, m)
	}
	out(w, "}")
	out(w, "")
}

func writeMethod(w io.Writer, m core.Method) {
	out(w, "    // %s", m.Description)
	out(w, `    func %s(input: %sInput) async throws -> %sOutput {`, m.LowerCamelName, m.CamelName, m.CamelName)
	out(w, `        return try await call(method: %q, input: input)`, m.Name)
	out(w, "    }")
	out(w, "")
}
