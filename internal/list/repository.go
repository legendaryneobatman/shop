package list

import (
	"fmt"
	"shop/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"shop/pkg/schema"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(userID int, input models.List) (int, error) {
	var list models.List
	query := fmt.Sprintf("INSERT INTO %s (user_id, title, description) VALUES ($1, $2, $3) RETURNING id, user_id, title, description", schema.ListTable)
	row := r.db.QueryRow(query, userID, input.Title, input.Description)
	err := row.Scan(&list.ID, &list.Title, &list.Description, &list.UserID)
	if err != nil {
		logrus.Errorf("error while creating list: %v", err)
		return 0, err
	}

	return list.ID, nil
}

func (r *Repository) GetAll(userID int) ([]models.List, error) {
	var lists []models.List
	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE user_id = $1", schema.ListTable)
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list models.List
		if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.UserID); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *Repository) GetByID(listID string) (models.List, error) {
	var list models.List
	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE id = $1", schema.ListTable)
	row := r.db.QueryRow(query, listID)
	if err := row.Scan(&list.ID, &list.UserID, &list.Title, &list.Description); err != nil {
		return models.List{}, err
	}
	return list, nil
}

func (r *Repository) Update(listID string, input models.List) (models.List, error) {
	var list models.List
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2 WHERE id = $3 RETURNING *", schema.ListTable)
	row := r.db.QueryRow(query, input.Title, input.Description, listID)
	if err := row.Scan(&list.ID, &list.UserID, &list.Title, &list.Description); err != nil {
		return models.List{}, err
	}

	return list, nil
}

func (r *Repository) GetWithPagination(userID int, limit int, offset int) ([]models.List, error) {
	var lists []models.List

	query := fmt.Sprintf(`
        SELECT id, title, description 
        FROM %s 
        WHERE user_id = $1 
        LIMIT $2 OFFSET $3`,
		schema.ListTable,
	)

	err := r.db.Select(&lists, query, userID, limit, offset)
	return lists, err
}

func (r *Repository) Delete(iDs []string) ([]string, error) {
	var deleteIds []string
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN($1) RETURNING id", schema.ListTable)
	rows, err := r.db.Query(query, iDs)
	if err != nil {
		logrus.Errorf("error while deleting lists: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		deleteIds = append(deleteIds, id)
	}
	return deleteIds, nil
}
