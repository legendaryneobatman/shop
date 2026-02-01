package models

type List struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	UserID      int    `json:"user_id" db:"user_id"`
}
