package product

import (
	"fmt"
	"shop/internal/models"

	"github.com/jmoiron/sqlx"
)

type ICategoryRepository interface {
	GetBySlug(slug string) (*models.Category, error)
	GetChildrenIDs(parentID int) ([]int, error)
}

type Repository struct {
	db                 *sqlx.DB
	categoryRepository ICategoryRepository
}

func NewRepository(db *sqlx.DB, categoryRepository ICategoryRepository) *Repository {
	return &Repository{
		db: db,
		categoryRepository: categoryRepository,
	}
}

func (r *Repository) GetProductsByCategoryPath(slug string) ([]*models.Product, error) {
	category, err := r.categoryRepository.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	childrenIDs, err := r.categoryRepository.GetChildrenIDs(category.ID)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT
     		p.id, p.sku, p.name, p.slug, p.description, p.short_description,
     		p.price, p.old_price, p.category_id, p.brand_id, p.quantity,
     		p.low_stock_threshold, p.main_image_url, p.weight, p.dimensions,
     		p.is_active, p.is_featured, p.is_new, p.rating, p.review_count,
     		p.created_at, p.updated_at, p.meta_title, p.meta_description,
     		p.meta_keywords
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE category_id in ($1)

	`)

	rows, err := r.db.Queryx(query, append(childrenIDs, category.ID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.StructScan(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
