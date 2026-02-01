package categories

import "go-shop/internal/models"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (r *Service) CreateCategory(input *models.Category) (*models.Category, error) {
	return &models.Category{}
}
