package auth

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go-shop/internal/user/entity"
	"go-shop/pkg/schema"
)

type ITokenRepository interface {
	SaveRefreshToken(token entity.RefreshToken) error
	GetRefreshTokenByHash(hash string) (*entity.RefreshToken, error)
	GetRefreshTokensByUserID(userID int) ([]*entity.RefreshToken, error)
	RevokeRefreshToken(tokenID string) error
	RevokeAllUserTokens(userID int) error
}

type RepositoryToken struct {
	db *sqlx.DB
}

func NewRepositoryToken(_db *sqlx.DB) *RepositoryToken {
	return &RepositoryToken{db: _db}
}

func (r *RepositoryToken) SaveRefreshToken(token entity.RefreshToken) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id,expires_at,ip_address,user_agent,revoked) values ($1, $2,$3, $4, $5)", schema.RefreshTokenTable)
	row := r.db.QueryRow(
		query,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.IPAddress,
		token.UserAgent,
		token.Revoked,
	)

	err := row.Scan()

	if err != nil {
		logrus.Errorf("Error when SaveRefreshToken %s", err.Error())
		return err
	}
	return nil
}
func (r *RepositoryToken) GetRefreshTokenByHash(hash string) (*entity.RefreshToken, error) {
	refreshToken := &entity.RefreshToken{}

	query := fmt.Sprintf(
		"SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, revoked FROM %s WHERE token_hash=$1",
		schema.RefreshTokenTable,
	)
	err := r.db.Get(refreshToken, query, hash)
	if err != nil {
		logrus.Errorf("Error when GetRefreshTokenByHash %s", err.Error())
		return nil, err
	}

	return refreshToken, nil
}
func (r *RepositoryToken) GetRefreshTokensByUserID(userID int) ([]entity.RefreshToken, error) {
	query := fmt.Sprintf(
		"SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, revoked FROM %s WHERE user_id=$1",
		schema.RefreshTokenTable,
	)
	rows, err := r.db.Queryx(query, userID)
	defer rows.Close()
	if rows.Err() != nil {
		logrus.Errorf("Error when rows for GetRefreshTokensByUserID %s", rows.Err().Error())
		return nil, rows.Err()
	}
	if err != nil {
		logrus.Errorf("Error when executing query for GetRefreshTokensByUserID %s", err.Error())
		return nil, err
	}

	var refreshTokens []entity.RefreshToken
	for rows.Next() {
		var rf entity.RefreshToken

		if err := rows.StructScan(&rf); err != nil {
			logrus.Fatalf("Error when scaning rows for GetRefreshTokensByUserID %s", err.Error())
			return nil, err
		}

		refreshTokens = append(refreshTokens, rf)
	}

	return refreshTokens, nil
}
func (r *RepositoryToken) RevokeRefreshToken(tokenID int) error {
	query := fmt.Sprintf("UPDATE %s SET revoked=true WHERE $1", schema.RefreshTokenTable)
	row := r.db.QueryRow(query, tokenID)
	err := row.Scan()
	if err != nil {
		logrus.Errorf("Error when RevokeRefreshToken %s", err.Error())
		return err
	}

	return nil
}
func (r *RepositoryToken) RevokeAllUserTokens(userID int) error {
	query := fmt.Sprintf("UPDATE %s SET revoked=true WHERE user_id=$1", schema.RefreshTokenTable)
	row := r.db.QueryRow(query, userID)
	err := row.Scan()

	if err != nil {
		logrus.Errorf("Error when RevokeAllUserTokens %s", err.Error())
	}

	return nil
}
