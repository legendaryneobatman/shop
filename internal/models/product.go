package models

import "time"

type Product struct {
	Id                int       `json:"id,omitempty" db:"id"`
	Sku               string    `json:"sku,omitempty" db:"sku"`
	Name              string    `json:"name,omitempty" db:"name"`
	Slug              string    `json:"slug,omitempty" db:"slug"`
	Description       string    `json:"description,omitempty" db:"description"`
	ShortDescription  string    `json:"shortDescription,omitempty" db:"short_description"`
	Price             float32   `json:"price,omitempty" db:"price"`
	OldPrice          float32   `json:"oldPrice,omitempty" db:"old_price"`
	CategoryId        int       `json:"categoryId,omitempty" db:"category_id"`
	BrandId           int       `json:"brandId,omitempty" db:"brand_id"`
	Quantity          int       `json:"quantity,omitempty" db:"quantity"`
	LowStockThreshold int       `json:"lowStockThreshold,omitempty" db:"low_stock_threshold"`
	MainImageUrl      string    `json:"mainImageUrl,omitempty" db:"main_image_url"`
	Weight            float32   `json:"weight,omitempty" db:"weight"`
	Dimensions        string    `json:"dimensions,omitempty" db:"dimensions"`
	IsActive          bool      `json:"isActive,omitempty" db:"is_active"`
	IsFeatured        bool      `json:"isFeatured,omitempty" db:"is_featured"`
	IsNew             bool      `json:"isNew,omitempty" db:"is_new"`
	Rating            float32   `json:"rating,omitempty" db:"rating"`
	ReviewCount       int       `json:"reviewCount,omitempty" db:"review_count"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
	MetaTitle         string    `json:"metaTitle,omitempty" db:"meta_title"`
	MetaDescription   string    `json:"metaDescription,omitempty" db:"meta_description"`
	MetaKeywords      string    `json:"metaKeywords,omitempty" db:"meta_keywords"`
}
