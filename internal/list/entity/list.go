package entity

type List struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	
	UserId int `json:"user_id"`
}
