package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-shop/internal/list/entity"
	"go-shop/pkg/schema"
)

type ListRepository struct {
	db *sqlx.DB
}

func NewListRepository(db *sqlx.DB) *ListRepository {
	return &ListRepository{db: db}
}

func (listRepository *ListRepository) Create(userId int, list entity.List) (int, error) {
	tx, err := listRepository.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	var user_id int
	var title string
	var description string
	createListQuery := fmt.Sprintf("INSERT INTO %s (user_id, title, description) VALUES ($1, $2, $3) RETURNING *", schema.TableNames.List)
	row := tx.QueryRow(createListQuery, userId, list.Title, list.Description)
	if err := row.Scan(&id, &user_id, &title, &description); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (listRepository *ListRepository) GetAll(userId int) ([]entity.List, error) {
	var lists []entity.List
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", schema.TableNames.List)
	rows, err := listRepository.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list entity.List
		if err := rows.Scan(&list.Id, &list.UserId, &list.Title, &list.Description); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func (listRepository *ListRepository) GetById(listId string) (entity.List, error) {
	var list entity.List
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", schema.TableNames.List)
	row := listRepository.db.QueryRow(query, listId)
	if err := row.Scan(&list.Id, &list.UserId, &list.Title, &list.Description); err != nil {
		return entity.List{}, err
	}
	return list, nil
}

func (listRepository *ListRepository) Update(listId string, input entity.List) (entity.List, error) {
	var list entity.List
	query := fmt.Sprintf("UPDATE %s SET title = $1, description = $2 WHERE id = $3 RETURNING *", schema.TableNames.List)
	row := listRepository.db.QueryRow(query, input.Title, input.Description, listId)
	if err := row.Scan(&list.Id, &list.UserId, &list.Title, &list.Description); err != nil {
		return entity.List{}, err
	}

	return list, nil
}

func (r *ListRepository) GetWithPagination(userId int, limit int, offset int) ([]entity.List, error) {
	var lists []entity.List

	query := fmt.Sprintf(`
        SELECT id, title, description 
        FROM %s 
        WHERE user_id = $1 
        ORDER BY id DESC 
        LIMIT $2 OFFSET $3`,
		schema.TableNames.List,
	)

	err := r.db.Select(&lists, query, userId, limit, offset)
	return lists, err
}
