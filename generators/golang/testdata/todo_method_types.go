type AddItemInput struct {
	// the item to add.
	Item todo.Item `json:"item"`
}

type AddItemOutput struct {
}

type GetItemsInput struct {
}

type GetItemsOutput struct {
	// Items is the list of to-do items.
	Items []todo.Item `json:"items"`
}

type RemoveItemInput struct {
	// the id of the item to remove.
	Id int `json:"id"`
}

type RemoveItemOutput struct {
	// the item removed.
	Item todo.Item `json:"item"`
}
