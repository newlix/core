// Item is a to-do item.
type Item struct {
	// ID is the unique id
	ID int `json:"id" db:"id"`

	// Text is the content
	Text string `json:"text" db:"text"`
}
