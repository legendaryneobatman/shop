package entity

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`

	ListID int `json:"list_id" binding:"required"`
}
