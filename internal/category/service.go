package category

import (
	"errors"
	"fmt"
	"go-shop/internal/models"

	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (r *Service) CreateCategory(input *models.Category) (*models.Category, error) {
	exists, err := r.repository.SlugExists(input.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, errors.New("category with this slug already exists")
	}

	if input.ParentID == nil || *input.ParentID == 0 {
		input.Level = 1
		input.Path = nil
	} else {
		parent, err := r.repository.GetByID(*input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent category does not exist")
		}

		if !parent.IsActive {
			return nil, errors.New("parent category is not active")
		}

		input.Level = parent.Level + 1

		input.Path = nil
	}

	category, err := r.repository.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	path, err := r.repository.BuildCategoryPath(category.Slug, category.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to build category path: %w", err)
	}

	category.Path = &path
	if err := r.repository.UpdatePath(category.ID, path); err != nil {
		return nil, fmt.Errorf("failed to update category path: %w", err)
	}

	return category, nil
}
func (r *Service) GetCategory(input *models.Category) (*models.Category, error) {
	category, err := r.repository.GetByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}

	return category, nil
}
func (r *Service) UpdateCategory(input *models.Category) (*models.Category, error) {
	category, err := r.repository.Update(input)
	if err != nil {
		logrus.Errorf("update category: %s", err.Error())
		return nil, err
	}
	return category, nil
}
func (r *Service) DeleteCategory(input *models.Category) (*models.Category, error) {

	return &models.Category{}, nil
}
func (r *Service) UpdateOrder(input *models.Category) (*models.Category, error) {
	return &models.Category{}, nil
}
func (r *Service) UpdateStatus(input *models.Category) (*models.Category, error) {
	return &models.Category{}, nil
}
func (r *Service) GetCategoryTree() ([]*models.CategoryNode, error) {
	//GET /api/category/tree — Самый важный эндпоинт.
	//Получение всего дерева категорий в виде вложенной структуры (Nested Set).
	//Ответ: Массив корневых категорий, каждая из которых содержит массив своих дочерних категорий (рекурсивно).
	//Только активные (is_active=true).
	//Для чего на фронте: Построение главного мега-меню, выпадающих списков в навигации, бокового меню фильтров.

	categories, err := r.repository.GetCategories()
	if err != nil {
		return nil, err
	}
	tree := buildTree(categories)

	return tree, nil
}
func (r *Service) GetCategories() ([]models.Category, error) {
	categories, err := r.repository.GetCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Service) GetCategoryBySlug(slug string) (*models.Category, error) {
	category, err := r.repository.GetBySlug(slug)
	if err != nil {return nil, err}
	return category, nil
}