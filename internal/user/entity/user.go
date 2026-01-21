package entity

import "time"

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`

	AvatarURL *string   `json:"avatar_url,omitempty" db:"avatar_url"` // Опциональное поле
	Phone     *string   `json:"phone,omitempty" db:"phone"`           // Опциональное поле
	Role      string    `json:"role" db:"role"`                       // Например: "user", "admin", "moderator"
	IsActive  bool      `json:"is_active" db:"is_active"`             // Для блокировки/активации
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
