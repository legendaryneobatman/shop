package models

import "time"

type Category struct {
	ID              int       `json:"id" db:"id"`
	ParentID        *int      `json:"parentId" db:"parent_id"`
	Name            string    `json:"name" db:"name"`
	Slug            string    `json:"slug" db:"slug"`
	Description     *string   `json:"description" db:"description"`
	ImageUrl        *string   `json:"imageUrl" db:"image_url"`
	DisplayOrder    int       `json:"displayOrder" db:"display_order"`
	IsActive        bool      `json:"isActive" db:"is_active"`
	MetaTitle       *string   `json:"metaTitle" db:"meta_title"`
	MetaDescription *string   `json:"metaDescription" db:"meta_description"`
	MetaKeywords    *string   `json:"metaKeywords" db:"meta_keywords"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
	Level           int       `json:"level" db:"level"`
	Path            *string   `json:"path" db:"path"`
}

type CategoryNode struct {
	Category

	Children []*CategoryNode `json:"children"`
}
