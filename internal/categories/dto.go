package categories

import (
	"go-shop/internal/models"
	"time"
)

type CreateRequestDTO struct {
	ParentID        int       `json:"parent_id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Slug            string    `json:"slug,omitempty"`
	Description     string    `json:"description,omitempty"`
	ImageUrl        string    `json:"image_url,omitempty"`
	DisplayOrder    int       `json:"display_order,omitempty"`
	IsActive        bool      `json:"is_active,omitempty"`
	MetaTitle       string    `json:"meta_title,omitempty"`
	MetaDescription string    `json:"meta_description,omitempty"`
	MetaKeywords    string    `json:"meta_keywords,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Level           int       `json:"level,omitempty"`
	Path            string    `json:"path,omitempty"`
}
type CreateResponseDTO = models.Category

func CategoryToCreateResponseDTO(input *CreateResponseDTO) *models.Category {
	return &models.Category{
		ID:              input.ID,
		ParentID:        input.ParentID,
		Name:            input.Name,
		Slug:            input.Slug,
		Description:     input.Description,
		ImageUrl:        input.ImageUrl,
		DisplayOrder:    input.DisplayOrder,
		IsActive:        input.IsActive,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
		MetaKeywords:    input.MetaKeywords,
		CreatedAt:       input.CreatedAt,
		UpdatedAt:       input.UpdatedAt,
		Level:           input.Level,
		Path:            input.Path,
	}
}
func CreateRequestDTOToCategory(input *CreateRequestDTO) *models.Category {
	return &models.Category{
		ParentID:        input.ParentID,
		Name:            input.Name,
		Slug:            input.Slug,
		Description:     input.Description,
		ImageUrl:        input.ImageUrl,
		DisplayOrder:    input.DisplayOrder,
		IsActive:        input.IsActive,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
		MetaKeywords:    input.MetaKeywords,
		CreatedAt:       input.CreatedAt,
		UpdatedAt:       input.UpdatedAt,
		Level:           input.Level,
		Path:            input.Path,
	}
}
