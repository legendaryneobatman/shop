package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	authService "go-shop/internal/auth/service"
	"net/http"
)

type SignInService struct {
	authService *authService.AuthService
}

func NewSignInService(authService *authService.AuthService) *SignInService {
	return &SignInService{
		authService: authService,
	}
}

func (p *SignInService) RenderSignInPage(c *gin.Context) {
	data := gin.H{}
	c.HTML(http.StatusOK, "sign-in", data)
}
func (p *SignInService) RetrySignIn(c *gin.Context) {
	c.Header("HX-Redirect", "/sign-in")
}
func (p *SignInService) HandleSignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	tokenPair, err := p.authService.Authenticate(username, password)
	if err != nil {
		logrus.Errorf("Error when verifying user %s", err.Error())
		c.HTML(http.StatusOK, "sign-in_wrong-credentials", gin.H{"Error": "Неверный логин или пароль"})
		return
	}

	c.SetCookie("Bearer", tokenPair.RefreshToken, 3600*24, "/", "", false, true)

	c.HTML(http.StatusOK, "auth-redirect", gin.H{
		"AccessToken": tokenPair.AccessToken,
	})
}