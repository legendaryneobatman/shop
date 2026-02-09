package list

import "go-shop/internal/models"

type IService interface {
	Create(userID int, list models.List) (int, error)
	GetAll(userID int) ([]models.List, error)
	GetWithPagination(userID int, limit int, offset int) ([]models.List, error)
	GetByID(listID string) (models.List, error)
	Update(listID string, input models.List) (models.List, error)
}

type Service struct {
	repo *Repository
}

func NewListService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID int, list models.List) (int, error) {
	return s.repo.Create(userID, list)
}
func (s *Service) GetAll(userID int) ([]models.List, error) {
	return s.repo.GetAll(userID)
}
func (s *Service) GetByID(listID string) (models.List, error) { return s.repo.GetByID(listID) }
func (s *Service) Update(listID string, input models.List) (models.List, error) {
	return s.repo.Update(listID, input)
}
func (s *Service) GetWithPagination(userID int, limit int, offset int) ([]models.List, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.GetWithPagination(userID, limit, offset)
}

func (s *Service) Delete(listIDs []string) ([]string, error) {
	return s.repo.Delete(listIDs)
}
