package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"shop/internal/models"
	"shop/internal/user"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims

	UserID int `json:"user_id"`
}

type Service struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	jwtSecretKey    []byte
	rUser           *user.Repository
	rToken          *RepositoryToken
}

func NewAuthService(rUser *user.Repository, rt *RepositoryToken) *Service {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		logrus.Fatalf("JWT_SECRET_KEY is not set in enviroment")
	}
	return &Service{
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		jwtSecretKey:    []byte(secret),
		rUser:           rUser,
		rToken:          rt,
	}
}

func (s *Service) Authenticate(username, password string) (*models.TokenPair, error) {
	_user, err := s.rUser.GetUserByUsername(username)
	if err != nil {
		logrus.Errorf("Error when trying to get _user for verify %s", err.Error())
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(_user.Password), []byte(password)); err != nil {
		logrus.Errorf("Error when compare hesh in Authenticate %s", err.Error())
		return nil, err
	}

	accessToken, err := s.generateAccessToken(_user.ID)
	if err != nil {
		logrus.Errorf("Error when try to generate access token in Authenticate %s", err.Error())
		return nil, err
	}

	refreshToken, _, err := s.generateRefreshToken()
	if err != nil {
		logrus.Errorf("Error when try to generate refresh token in Authenticate %s", err.Error())
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *Service) ParseTokenForUserID(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtSecretKey, nil
	})

	if err != nil {
		logrus.Errorf("Error when try to parse token %s", err.Error())
		return 0, err
	}

	if !token.Valid {
		logrus.Errorf("Token is not valid")
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		logrus.Errorf("Token claims are not of type *tokenClaims")
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}
func (s *Service) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	tHash := hex.EncodeToString(hash[:])

	storedRefreshT, err := s.rToken.GetRefreshTokenByHash(tHash)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	storedTExpiresAtTime, err := time.Parse("YYYY-MM-DD HH:MM:SS", storedRefreshT.ExpiresAt)
	if err != nil {
		logrus.Errorf("stored refresh token is expired %s", err.Error())
		return nil, err
	}

	if time.Now().After(storedTExpiresAtTime) {
		logrus.Errorf("stored refresh token is expired")
		return nil, err
	}

	if storedRefreshT.Revoked != nil && *storedRefreshT.Revoked {
		logrus.Errorf("Stored token is revoked")
		return nil, err
	}

	_user, err := s.rUser.GetUserByID(storedRefreshT.UserID)
	if err != nil {
		logrus.Errorf("Error when geting _user by id %s", err.Error())
		return nil, err
	}

	newAccessToken, err := s.generateAccessToken(_user.ID)
	if err != nil {
		logrus.Errorf("Failed to generate access token %s", err.Error())
		return nil, err
	}

	if err := s.rToken.RevokeRefreshToken(storedRefreshT.ID); err != nil {
		logrus.Errorf("Failed to revoke refresh token %s", err.Error())
		return nil, err
	}

	newRefreshToken, newRefreshTokenHash, err := s.generateRefreshToken()
	if err != nil {
		logrus.Errorf("Failed to generate refresh token %s", err.Error())
		return nil, err
	}
	newTokenEntity := models.RefreshToken{
		UserID:    _user.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(s.RefreshTokenTTL).String(),
		Revoked:   nil,
	}
	if err := s.rToken.SaveRefreshToken(newTokenEntity); err != nil {
		logrus.Errorf("Error when saving refresh token %s", err.Error())
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
func (s *Service) RevokeToken(refreshToken string) error {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	storedToken, err := s.rToken.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		logrus.Errorf("Error when searching for stored token by hash %s", err.Error())
		return err
	}

	if err := s.rToken.RevokeRefreshToken(storedToken.ID); err != nil {
		logrus.Errorf("Error when revoking token %s", err.Error())
		return err
	}

	return nil
}
func (s *Service) RevokeAllTokens(userID int) error {
	if err := s.rToken.RevokeAllUserTokens(userID); err != nil {
		logrus.Errorf("Error when revoking all tokens %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) generateRefreshToken() (token string, tokenHash string, err error) {
	token = uuid.New().String()

	hash := sha256.Sum256([]byte(token))
	tokenHash = hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}
func (s *Service) generateAccessToken(userID int) (string, error) {
	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(userID),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecretKey)
}
