package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-shop/internal/user/entity"
	"go-shop/pkg/schema"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", schema.TableNames.User)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", schema.TableNames.User)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
