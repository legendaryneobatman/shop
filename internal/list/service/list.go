package service

import (
	"go-shop/internal/list/entity"
	"go-shop/internal/list/repository"
)

type ListService struct {
	repo *repository.ListRepository
}

func NewListService(repo *repository.ListRepository) *ListService {
	return &ListService{repo: repo}
}

func (s *ListService) Create(userId int, list entity.List) (int, error) {
	return s.repo.Create(userId, list)
}
func (s *ListService) GetAll(userId int) ([]entity.List, error) {
	return s.repo.GetAll(userId)
}
func (s *ListService) GetById(listId string) (entity.List, error) { return s.repo.GetById(listId) }
func (s *ListService) Update(listId string, input entity.List) (entity.List, error) {
	return s.repo.Update(listId, input)
}
func (s *ListService) GetWithPagination(userId int, limit int, offset int) ([]entity.List, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.GetWithPagination(userId, limit, offset)
}
