package user

import (
	"shop/internal/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(input *models.User) (models.User, error) {
	input.Password = string(s.generatePasswordHash(input.Password))
	_user := models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	}
	id, err := s.repo.CreateUser(&_user)
	if err != nil {
		logrus.Errorf("Error when CreateUser %s", err.Error())
		return _user, err
	}
	_user.ID = id
	return _user, nil
}

func (s *Service) generatePasswordHash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to encrypt password error: %s", err.Error())
	}

	return hash
}
