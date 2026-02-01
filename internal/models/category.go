package models

import "time"

type Category struct {
	ID              int       `json:"id" json:"id"`
	ParentID        int       `json:"parent_id" json:"parentId"`
	Name            string    `json:"name" json:"name"`
	Slug            string    `json:"slug" json:"slug"`
	Description     string    `json:"description" json:"description"`
	ImageUrl        string    `json:"image_url" json:"imageUrl"`
	DisplayOrder    int       `json:"display_order" json:"displayOrder"`
	IsActive        bool      `json:"is_active" json:"isActive"`
	MetaTitle       string    `json:"meta_title" json:"metaTitle"`
	MetaDescription string    `json:"meta_description" json:"metaDescription"`
	MetaKeywords    string    `json:"meta_keywords" json:"metaKeywords"`
	CreatedAt       time.Time `json:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `json:"updated_at" json:"updatedAt"`
	Level           int       `json:"level" json:"level"`
	Path            string    `json:"path" json:"path"`
}
