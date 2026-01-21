package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth/repository"
	"go-shop/internal/user/entity"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

const (
	salt         = "aslkdjaslk"
	signingUpKey = "aksldjhaksjhdakjshd"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	jwtSecretKey    []byte
	ra              *repository.AuthRepository
	rt              *repository.TokenRepository
}

func NewAuthService(ra *repository.AuthRepository, rt *repository.TokenRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		logrus.Fatalf("JWT_SECRET_KEY is not set in enviroment")
	}
	return &AuthService{
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 7 * 24 * time.Hour,
		ra:              ra,
		rt:              rt,
		jwtSecretKey:    []byte(secret),
	}
}

func (s *AuthService) CreateUser(user *entity.User) (int, error) {
	user.Password = string(s.generatePasswordHash(user.Password))
	return s.ra.CreateUser(user)
}
func (s *AuthService) Authenticate(username, password string) (*TokenPair, error) {
	user, err := s.ra.GetUserByUsername(username)
	if err != nil {
		logrus.Errorf("Error when trying to get user for verify %s", err.Error())
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.Errorf("Error when compare hesh in Authenticate %s", err.Error())
		return nil, err
	}

	accessToken, err := s.generateAccessToken(user.Id)
	if err != nil {
		logrus.Errorf("Error when try to generate access token in Authenticate %s", err.Error())
		return nil, err
	}

	refreshToken, _, err := s.generateRefreshToken()
	if err != nil {
		logrus.Errorf("Error when try to generate refresh token in Authenticate %s", err.Error())
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.rt.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		logrus.Errorf("Error when searching for stored token by hash %s", err.Error())
		return nil, err
	}

	storedTokenExpiresAtTime, err := time.Parse("YYYY-MM-DD HH:MM:SS", storedToken.ExpiresAt)
	if err != nil {
		logrus.Errorf("Error when parsing time from storedToken %s", err.Error())
		return nil, err
	}

	if time.Now().After(storedTokenExpiresAtTime) {
		logrus.Errorf("Stored token is expired")
		return nil, err
	}

	if storedToken.Revoked != nil && *storedToken.Revoked {
		logrus.Errorf("Stored token is revoked")
		return nil, err
	}

	user, err := s.ra.GetUserById(storedToken.UserId)
	if err != nil {
		logrus.Errorf("Error when geting user by id %s", err.Error())
		return nil, err
	}

	newAccessToken, err := s.generateAccessToken(user.Id)
	if err != nil {
		logrus.Errorf("Failed to generate access token %s", err.Error())
		return nil, err
	}

	if err := s.rt.RevokeRefreshToken(storedToken.Id); err != nil {
		logrus.Errorf("Failed to revoke refresh token %s", err.Error())
		return nil, err
	}

	newRefreshToken, newRefreshTokenHash, err := s.generateRefreshToken()
	if err != nil {
		logrus.Errorf("Failed to generate refresh token %s", err.Error())
		return nil, err
	}

	newTokenEntity := entity.RefreshToken{
		UserId:    user.Id,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL).String(),
		Revoked:   nil,
	}

	if err := s.rt.SaveRefreshToken(newTokenEntity); err != nil {
		logrus.Errorf("Error when saving refresh token %s", err.Error())
		return nil, nil
	}

	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
func (s *AuthService) RevokeToken(refreshToken string) error {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.rt.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		logrus.Errorf("Error when searching for stored token by hash %s", err.Error())
		return err
	}

	if err := s.rt.RevokeRefreshToken(storedToken.Id); err != nil {
		logrus.Errorf("Error when revoking token %s", err.Error())
		return err
	}

	return nil
}
func (s *AuthService) RevokeAllTokens(userId int) error {
	if err := s.rt.RevokeAllUserTokens(userId); err != nil {
		logrus.Errorf("Error when revoking all tokens %s", err.Error())
		return err
	}

	return nil
}

func (s *AuthService) generatePasswordHash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed to encrypt password error: %s", err.Error())
	}

	return hash
}
func (s *AuthService) generateRefreshToken() (token string, tokenHash string, err error) {
	token = uuid.New().String()

	hash := sha256.Sum256([]byte(token))
	tokenHash = hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}
func (s *AuthService) generateAccessToken(userId int) (string, error) {
	claims := &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userId),
		},
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecretKey)
}
