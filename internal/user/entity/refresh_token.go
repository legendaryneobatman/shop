package entity

type RefreshToken struct {
	Id        int     `json:"id" db:"id"`
	UserId    int     `json:"user_id" db:"user_id"`
	TokenHash string  `json:"token_hash" db:"token_hash"`
	ExpiresAt string  `json:"expires_at" db:"expires_at"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
	Revoked   *bool   `json:"revoked" db:"revoked"`
}
