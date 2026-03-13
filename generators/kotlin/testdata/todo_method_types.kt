@Serializable
data class AddItemInput(
    @SerialName("item") val item: Item = Item(),
)

@Serializable
class AddItemOutput(
)

@Serializable
class GetItemsInput(
)

@Serializable
data class GetItemsOutput(
    @SerialName("items") val items: List<Item> = emptyList(),
)

@Serializable
data class RemoveItemInput(
    @SerialName("id") val id: Long = 0,
)

@Serializable
data class RemoveItemOutput(
    @SerialName("item") val item: Item = Item(),
)

