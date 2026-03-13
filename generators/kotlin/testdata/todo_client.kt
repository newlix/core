data class CoreError(
    val status: Int,
    val msg: String
) : Exception()

// TodoClient is the API client.
// url is the required API endpoint address.
class TodoClient(val endpoint: String) {
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

    // AddItem adds an item to the list.
    suspend fun addItem(input: AddItemInput): AddItemOutput {
        return call("add_item", input)
    }

    // GetItems returns all items in the list.
    suspend fun getItems(input: GetItemsInput): GetItemsOutput {
        return call("get_items", input)
    }

    // RemoveItem removes an item from the to-do list.
    suspend fun removeItem(input: RemoveItemInput): RemoveItemOutput {
        return call("remove_item", input)
    }

}

