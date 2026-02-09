package category

import (
	"go-shop/internal/models"
	"time"
)

type CreateRequestDTO struct {
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
}
type CreateResponseDTO struct {
	ID              int       `json:"id"`
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
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
func ToCreateResponseDTO(input *models.Category) *CreateResponseDTO {
	return &CreateResponseDTO{
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

type GetResponseDTO struct {
	ID              int       `json:"id"`
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
}

func ToGetResponseDTO(input *models.Category) *GetResponseDTO {
	return &GetResponseDTO{
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

type UpdateRequestDTO struct {
	ParentID        *int   `json:"parent_id,omitempty" binding:"required"`
	Name            string `json:"name,omitempty" binding:"required"`
	Slug            string `json:"slug,omitempty" binding:"required"`
	Description     string `json:"description,omitempty" binding:"required"`
	ImageUrl        string `json:"image_url,omitempty" binding:"required"`
	DisplayOrder    int    `json:"display_order,omitempty" binding:"required"`
	IsActive        bool   `json:"is_active,omitempty" binding:"required"`
	MetaTitle       string `json:"meta_title,omitempty" binding:"required"`
	MetaDescription string `json:"meta_description,omitempty" binding:"required"`
	MetaKeywords    string `json:"meta_keywords,omitempty" binding:"required"`
	Level           int    `json:"level,omitempty" binding:"required"`
}
type UpdateResponseDTO struct {
	ID              int       `json:"id"`
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
}

func UpdateRequestDTOToCategory(ID int, input *UpdateRequestDTO) *models.Category {
	return &models.Category{
		ID:              ID,
		ParentID:        input.ParentID,
		Name:            input.Name,
		Slug:            input.Slug,
		Description:     &input.Description,
		ImageUrl:        &input.ImageUrl,
		DisplayOrder:    input.DisplayOrder,
		IsActive:        input.IsActive,
		MetaTitle:       &input.MetaTitle,
		MetaDescription: &input.MetaDescription,
		MetaKeywords:    &input.MetaKeywords,
		Level:           input.Level,
	}

}
func ToUpdateResponseDTO(input *models.Category) *UpdateResponseDTO {
	return &UpdateResponseDTO{
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

type DeleteResponseDTO struct {
	ID int `json:"id,omitempty"`
}

func DeleteRequestDTOToCategory(ID int) *models.Category {
	return &models.Category{
		ID: ID,
	}
}
func ToDeleteResponseDTO(input *models.Category) *DeleteResponseDTO {
	return &DeleteResponseDTO{
		ID: input.ID,
	}
}

type UpdateOrderRequestDTO struct {
	Level int `json:"level"`
}
type UpdateOrderResponseDTO struct {
	ID              int       `json:"id"`
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
}

func UpdateOrderRequestDTOToCategory(ID int, input *UpdateOrderRequestDTO) *models.Category {
	return &models.Category{
		ID:    ID,
		Level: input.Level,
	}
}
func ToUpdateOrderResponseDTO(input *models.Category) *UpdateOrderResponseDTO {
	return &UpdateOrderResponseDTO{
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

type UpdateStatusRequestDTO struct {
	IsActive bool `json:"isActive"`
}
type UpdateStatusResponseDTO struct {
	ID              int       `json:"id"`
	ParentID        *int      `json:"parentId"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description"`
	ImageUrl        *string   `json:"imageUrl"`
	DisplayOrder    int       `json:"displayOrder"`
	IsActive        bool      `json:"isActive"`
	MetaTitle       *string   `json:"metaTitle"`
	MetaDescription *string   `json:"metaDescription"`
	MetaKeywords    *string   `json:"metaKeywords"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Level           int       `json:"level"`
	Path            *string   `json:"path"`
}

func UpdateStatusRequestDTOToCategory(ID int, input *UpdateStatusRequestDTO) *models.Category {
	return &models.Category{
		ID:       ID,
		IsActive: input.IsActive,
	}
}
func ToUpdateStatusResponseDTO(input *models.Category) *UpdateStatusResponseDTO {
	return &UpdateStatusResponseDTO{
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
