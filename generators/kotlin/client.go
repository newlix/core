package kotlin

import (
	"io"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateClientFileConfig struct {
	Output       string
	Package      string
	Methods      []core.Method
	TypesPackage string
	Client       string
}

func GenerateClientFile(c GenerateClientFileConfig) error {
	return common.GenerateFile(c.Output, func(w io.Writer) error {
		common.GenerateWarning(w)
		out(w, "@file:Suppress(\"unused\")")
		out(w, "package %s", c.Package)
		out(w, "")
		out(w, "import kotlinx.coroutines.Dispatchers")
		out(w, "import kotlinx.coroutines.withContext")
		out(w, "import kotlinx.serialization.encodeToString")
		out(w, "import kotlinx.serialization.json.Json")
		out(w, "import kotlinx.serialization.SerialName")
		out(w, "import kotlinx.serialization.Serializable")
		out(w, "import okhttp3.OkHttpClient")
		out(w, "import okhttp3.Request")
		out(w, "import okhttp3.RequestBody.Companion.toRequestBody")
		out(w, "")
		if c.TypesPackage != "" {
			out(w, "import %s.*", c.TypesPackage)
		}
		GenerateClient(w, c.Methods, c.Client)
		GenerateMethodTypes(w, c.Methods)
		return nil
	})
}

const start = `data class CoreError(
    val status: Int,
    val msg: String
) : Exception("HTTP $status: $msg")

// %s is the API client.
// url is the required API endpoint address.
class %s(val endpoint: String, val client: OkHttpClient = OkHttpClient()) {
    private val json = Json { ignoreUnknownKeys = true }

    // AuthToken is an optional authentication token.
    var authToken: String? = null

    private suspend inline fun <reified Input, reified Output> call(
        method: String, input: Input
    ): Output {
        return withContext(Dispatchers.IO) {
            val url = "$endpoint/$method"
            val postBody = json.encodeToString(input)
            val request = Request.Builder()
                .url(url)
                .post(postBody.toRequestBody())
                .addHeader("Content-Type", "application/json")
            if (authToken != null) {
                request.addHeader("Authorization", "Bearer $authToken")
            }

            return@withContext client.newCall(request.build()).execute().use { response ->
                val body: String = response.body?.string()
                    ?: throw CoreError(status = response.code, msg = "empty response body")
                if (!response.isSuccessful) {
                    throw CoreError(
                        status = response.code,
                        msg = body
                    )
                }
                return@use json.decodeFromString<Output>(body)
            }
        }
    }
`

func GenerateClient(w io.Writer, mm []core.Method, client string) {
	out(w, start, client, client)
	writeMethods(w, mm)
	out(w, "}")
	out(w, "")
}

func writeMethods(w io.Writer, mm []core.Method) {
	for _, m := range mm {
		out(w, "    // %s", m.Description)
		out(w, "    suspend fun %s(input: %sInput): %sOutput {", m.LowerCamelName, m.CamelName, m.CamelName)
		out(w, `        return call(%q, input)`, m.Name)
		out(w, "    }")
		out(w, "")
	}
}
