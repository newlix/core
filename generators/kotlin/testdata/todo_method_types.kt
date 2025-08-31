@Serializable
data class AddItemInput(
    @SerialName("item") var item: Item = Item(),
)

@Serializable
data class AddItemOutput(
)

@Serializable
data class GetItemsInput(
)

@Serializable
data class GetItemsOutput(
    @SerialName("items") var items: MutableList<Item> = mutableListOf(),
)

@Serializable
data class RemoveItemInput(
    @SerialName("id") var id: Int = 0,
)

@Serializable
data class RemoveItemOutput(
    @SerialName("item") var item: Item = Item(),
)

