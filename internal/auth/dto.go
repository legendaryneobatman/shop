package auth

type SignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type SignUpOutput struct {
	ID int `json:"id"`
}

type SignInInput struct {
	Username string `json:"username" binding:"required" json:"username"`
	Password string `json:"password" binding:"required" json:"password"`
}
type SignInOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
type RefreshOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
