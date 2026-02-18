package user

import (
	"fmt"
	"shop/internal/models"
	"shop/pkg/schema"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", schema.UserTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUser(username, password string) (*models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM %s WHERE username=$1 AND password_hash=$2", schema.UserTable)
	err := r.db.Get(&user, query, username, password)

	return &user, err
}

func (r *Repository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM %s WHERE id=$1",
		schema.UserTable,
	)
	err := r.db.Get(&user, query, userID)
	if err != nil {
		logrus.Errorf("Error when GetUserByID %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM %s WHERE username=$1",
		schema.UserTable,
	)
	err := r.db.Get(&user, query, username)
	if err != nil {
		logrus.Errorf("Error when GetUserByUsername %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetAll() ([]models.User, error) {
	users := []models.User{}
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, avatar_url, email, phone, role, is_active, created_at, updated_at FROM %s",
		schema.UserTable,
	)
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}
