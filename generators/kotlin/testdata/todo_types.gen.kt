// Item is a to-do item.
@Serializable
data class Item(
    @SerialName("id") val id: Long = 0,
    @SerialName("text") val text: String = "",
    @SerialName("created_at") val createdAt: Date = Date(),
)

