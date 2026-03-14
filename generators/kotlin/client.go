package kotlin

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateClientFileConfig struct {
	Output       string
	Package      string
	Methods      []core.Method
	TypesPackage string
	Types        []core.Type
	Client       string
}

func GenerateClientFile(c GenerateClientFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	common.GenerateWarning(w)
	out(w, "@file:Suppress(\"unused\")")
	out(w, "package %s", c.Package)
	out(w, `
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody

`)
	if c.TypesPackage != "" {
		out(w, "import %s.*", c.TypesPackage)
	}
	GenerateClient(w, c.Methods, c.Client)
	GenerateMethodTypes(w, c.Methods)
}

var start = `data class CoreError(
    val status: Int,
    val msg: String
) : Exception()

// %s is the API client.
// url is the required API endpoint address.
class %s(val endpoint: String) {
    private val json = Json { ignoreUnknownKeys = true }

    // AuthToken is an optional authentication token.
    var authToken: String? = null

    // client is used for making requests.
    val client = OkHttpClient()

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
                val body: String = response.body!!.string()
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
