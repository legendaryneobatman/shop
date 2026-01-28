package list

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-shop/pkg/schema"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (lr *Repository) Create(userID int, list List) (int, error) {
	tx, err := lr.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	var _userID int
	var title string
	var description string
	createListQuery := fmt.Sprintf("INSERT INTO %s (_userID, title, description) VALUES ($1, $2, $3) RETURNING *", schema.ListTable)
	row := tx.QueryRow(createListQuery, userID, list.Title, list.Description)
	if err := row.Scan(&id, &title, &description, &_userID); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	return id, tx.Commit()
}

func (lr *Repository) GetAll(userID int) ([]List, error) {
	var lists []List
	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE user_id = $1", schema.ListTable)
	rows, err := lr.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list List
		if err := rows.Scan(&list.ID, &list.UserID, &list.Title, &list.Description); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func (lr *Repository) GetByID(listID string) (List, error) {
	var list List
	query := fmt.Sprintf("SELECT id, title, description, user_id FROM %s WHERE id = $1", schema.ListTable)
	row := lr.db.QueryRow(query, listID)
	if err := row.Scan(&list.ID, &list.UserID, &list.Title, &list.Description); err != nil {
		return List{}, err
	}
	return list, nil
}

func (lr *Repository) Update(listID string, input List) (List, error) {
	var list List
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2 WHERE id = $3 RETURNING *", schema.ListTable)
	row := lr.db.QueryRow(query, input.Title, input.Description, listID)
	if err := row.Scan(&list.ID, &list.UserID, &list.Title, &list.Description); err != nil {
		return List{}, err
	}

	return list, nil
}

func (lr *Repository) GetWithPagination(userID int, limit int, offset int) ([]List, error) {
	var lists []List

	query := fmt.Sprintf(`
        SELECT id, title, description 
        FROM %s 
        WHERE user_id = $1 
        ORDER BY id DESC 
        LIMIT $2 OFFSET $3`,
		schema.ListTable,
	)

	err := lr.db.Select(&lists, query, userID, limit, offset)
	return lists, err
}
