package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth/repository"
	"go-shop/internal/user/entity"
	"time"
)

const (
	salt         = "aslkdjaslk"
	signingUpKey = "aksldjhaksjhdakjshd"
	tokenTTL     = 12 * time.Hour
)

type AuthService struct {
	repository *repository.AuthRepository
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repository *repository.AuthRepository) *AuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repository.CreateUser(user)
}

func (s *AuthService) VerifyUser(username, password string) (string, error) {
	user, err := s.repository.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		logrus.Errorf("Error when trying to get user for verify %s", err.Error())

		return "", err

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(signingUpKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingUpKey), nil
	})

	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
