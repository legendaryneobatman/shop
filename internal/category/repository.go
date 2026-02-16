package category

import (
	"database/sql"
	"errors"
	"fmt"
	"go-shop/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(input *models.Category) (*models.Category, error) {
	query := `
        INSERT INTO categories (
            parent_id, name, slug, description, image_url,
            display_order, is_active, meta_title, meta_description,
            meta_keywords, level, path
        )
        VALUES (
            :parent_id, :name, :slug, :description, :image_url,
            :display_order, :is_active, :meta_title, :meta_description,
            :meta_keywords, :level, :path
        )
        RETURNING id, parent_id, name, slug, description, image_url,
                  display_order, is_active, meta_title, meta_description,
                  meta_keywords, created_at, updated_at, level, path
    `

	rows, err := r.db.NamedQuery(query, input)
	if err != nil {
		logrus.Errorf("Error when executing query for category create %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var result models.Category
	if rows.Next() {
		err = rows.StructScan(&result)
		if err != nil {
			logrus.Errorf("Error when filling struct for category create %s", err.Error())
			return nil, err
		}
	}

	return &result, nil
}

func (r *Repository) GetCategories() ([]models.Category, error) {
	query := fmt.Sprintf("SELECT * FROM categories")
	categories := []models.Category{}
	err := r.db.Select(&categories, query)
	if err != nil {
		logrus.Errorf("Failed to select categories %v", err.Error())
		return nil, err
	}
	return categories, nil
}
func (r *Repository) SlugExists(slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1)`

	var exists bool
	err := r.db.Get(&exists, query, slug)
	if err != nil {
		logrus.Errorf("Error checking slug existence: %s", err.Error())
		return false, err
	}

	return exists, nil
}
func (r *Repository) CheckSlugsExist(slugs []string) (map[string]bool, error) {
	if len(slugs) == 0 {
		return make(map[string]bool), nil
	}

	query := `SELECT slug FROM categories WHERE slug = ANY($1)`

	var existingSlugs []string
	err := r.db.Select(&existingSlugs, query, pq.Array(slugs))
	if err != nil {
		return nil, fmt.Errorf("failed to check slugs: %w", err)
	}

	result := make(map[string]bool, len(existingSlugs))
	for _, slug := range existingSlugs {
		result[slug] = true
	}

	return result, nil
}
func (r *Repository) GetByID(id int) (*models.Category, error) {
	query := `
        SELECT id, parent_id, name, slug, description, image_url,
               display_order, is_active, meta_title, meta_description,
               meta_keywords, created_at, updated_at, level, path
        FROM categories
        WHERE id = $1
    `

	var category models.Category
	err := r.db.Get(&category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.Errorf("Error getting category by ID: %s", err.Error())
		return nil, err
	}

	return &category, nil
}
func (r *Repository) UpdatePath(id int, path string) error {
	query := `UPDATE categories SET path = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(query, path, id)
	if err != nil {
		logrus.Errorf("Error updating category path: %s", err.Error())
		return err
	}

	return nil
}
func (r *Repository) GetChildrenIDs(parentID int) ([]int, error) {
	query := `
        WITH RECURSIVE category_tree AS (
            -- Базовый случай: прямые потомки
            SELECT id FROM categories WHERE parent_id = $1
            UNION ALL
            -- Рекурсивный случай: потомки потомков
            SELECT c.id 
            FROM categories c
            INNER JOIN category_tree ct ON c.parent_id = ct.id
        )
        SELECT id FROM category_tree
    `

	var ids []int
	err := r.db.Select(&ids, query, parentID)
	if err != nil {
		logrus.Errorf("Error getting children IDs: %s", err.Error())
		return nil, err
	}

	return ids, nil
}
func (r *Repository) UpdateLevelAndPath(id int, level int, path string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE categories SET level = $1, path = $2, updated_at = NOW() WHERE id = $3`
	_, err = tx.Exec(query, level, path, id)
	if err != nil {
		logrus.Errorf("Error updating category level and path: %s", err.Error())
		return err
	}

	var children []models.Category
	childQuery := `
        SELECT id, parent_id, name, slug, description, image_url,
               display_order, is_active, meta_title, meta_description,
               meta_keywords, created_at, updated_at, level, path
        FROM categories 
        WHERE parent_id = $1
    `
	err = tx.Select(&children, childQuery, id)
	if err != nil {
		logrus.Errorf("Error getting children: %s", err.Error())
		return err
	}

	for _, child := range children {
		childLevel := level + 1
		childPath := fmt.Sprintf("%s%d/", path, child.ID)

		if err := r.updateLevelAndPathRecursive(tx, child.ID, childLevel, childPath); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) updateLevelAndPathRecursive(tx *sqlx.Tx, id int, level int, path string) error {
	query := `UPDATE categories SET level = $1, path = $2, updated_at = NOW() WHERE id = $3`
	_, err := tx.Exec(query, level, path, id)
	if err != nil {
		return err
	}

	var children []models.Category
	childQuery := `
        SELECT id, parent_id, name, slug, description, image_url,
               display_order, is_active, meta_title, meta_description,
               meta_keywords, created_at, updated_at, level, path
        FROM categories 
        WHERE parent_id = $1
    `
	err = tx.Select(&children, childQuery, id)
	if err != nil {
		return err
	}

	for _, child := range children {
		childLevel := level + 1
		childPath := fmt.Sprintf("%s%d/", path, child.ID)

		if err := r.updateLevelAndPathRecursive(tx, child.ID, childLevel, childPath); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Update(input *models.Category) (*models.Category, error) {
	query := fmt.Sprintf(
		`
				UPDATE categories
				SET 
				    parent_id= :parent_id,
    				name= :name,
				    slug= :slug,
    				description= :description,
    				image_url= :image_url,
    				display_order= :display_order,
    				is_active= :is_active,
    				meta_title= :meta_title,
    				meta_description= :meta_description,
    				meta_keywords= :meta_keywords,
    				updated_at= :updated_at
    			WHERE id= :id
				RETURNING id, parent_id, name, slug, description, image_url,
				          display_order, is_active, meta_title, meta_description,
				          meta_keywords, created_at, updated_at, level, path
		`,
	)
	rows, err := r.db.NamedQuery(query, input)
	if err != nil {
		logrus.Errorf("Error when updating category: %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	var category models.Category
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err := rows.StructScan(&category); err != nil {
		logrus.Errorf("Error scanning updated category: %s", err.Error())
		return nil, err
	}
	path, err := r.BuildCategoryPath(category.Slug, category.ParentID)
	if err != nil {
		logrus.Errorf("Cant rebuild path when update: %s", err.Error())
		return nil, err
	}
	category.Path = &path
	if err := r.UpdatePath(category.ID, path); err != nil {
		return nil, fmt.Errorf("failed to update category path: %w", err)
	}

	return &category, nil
}

func (r *Repository) BuildCategoryPath(slug string, parentID *int) (string, error) {
	if parentID == nil || *parentID == 0 {
		return fmt.Sprintf("/%s/", slug), nil
	}

	parent, err := r.GetByID(*parentID)
	if err != nil {
		return "", err
	}
	if parent == nil {
		return "", errors.New("parent category not found")
	}

	if parent.Path != nil {
		return fmt.Sprintf("%s%s/", *parent.Path, slug), nil
	}

	return fmt.Sprintf("/%s/%s/", parent.Slug, slug), nil
}

func (r *Repository) GetBySlug(slug string) (*models.Category, error) {
	query := `
        SELECT id, parent_id, name, slug, description, image_url,
               display_order, is_active, meta_title, meta_description,
               meta_keywords, created_at, updated_at, level, path
        FROM categories
        WHERE slug = $1
    `

	var category models.Category
	err := r.db.Get(&category, query, slug)
	if err != nil {
		logrus.Errorf("Error getting category by slug: %s", err.Error())
		return nil, err
	}

	return &category, nil
}

func (r *Repository) Delete(ID int) (int, error) {
	var _ID int
	query := fmt.Sprintf("DELETE FROM categories WHERE id=$1 RETURNING id")
	err := r.db.QueryRow(query, ID).Scan(&_ID)

	if err != nil {
		logrus.Errorf("error when scan row for deleted category: %v", err)
		return 0, err
	}

	return _ID, nil
}
