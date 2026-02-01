package models

import "time"

type Brand struct {
	Id              int       `json:"id,omitempty" db:"id"`
	Name            string    `json:"name,omitempty" db:"name"`
	Slug            string    `json:"slug,omitempty" db:"slug"`
	Description     string    `json:"description,omitempty" db:"description"`
	LogoUrl         string    `json:"logo_url,omitempty" db:"logo_url"`
	WebsiteUrl      string    `json:"website_url,omitempty" db:"website_url"`
	IsActive        bool      `json:"is_active,omitempty" db:"is_active"`
	IsFeatured      bool      `json:"is_featured,omitempty" db:"is_featured"`
	DisplayOrder    int       `json:"display_order,omitempty" db:"display_order"`
	MetaTitle       string    `json:"meta_title,omitempty" db:"meta_title"`
	MetaDescription string    `json:"meta_description,omitempty" db:"meta_description"`
	MetaKeywords    string    `json:"meta_keywords,omitempty" db:"meta_keywords"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
