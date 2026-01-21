package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go-shop/internal/user/entity"
	"go-shop/pkg/schema"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user *entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", schema.TableNames.User)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (*entity.User, error) {
	var user entity.User

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", schema.TableNames.User)
	err := r.db.Get(&user, query, username, password)

	return &user, err
}

func (r *AuthRepository) GetUserById(userId int) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", schema.TableNames.User)
	err := r.db.Get(&user, query, userId)
	if err != nil {

	}

	if err != nil {
		logrus.Errorf("Error when GetUserById %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", schema.TableNames.User)
	err := r.db.Get(&user, query, username)
	if err != nil {
		logrus.Errorf("Error when GetUserByUsername %s", err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetAll() (*[]entity.User, error) {
	var users []entity.User
	query := fmt.Sprintf("SELECT * FROM %s", schema.TableNames.List)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}
