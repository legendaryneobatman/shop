package entity

type List struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
