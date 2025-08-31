// Item is a to-do item.
@Serializable
data class Item(
    @SerialName("id") var id: Int = 0,
    @SerialName("text") var text: String = "",
    @SerialName("created_at") var createdAt: Date = Date(),
)

