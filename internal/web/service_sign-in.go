package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth"
	"net/http"
)

const (
	cookieMaxAge = 3600 * 24 // 24 hours
)

type ServiceSignIn struct {
	authService *auth.Service
}

func NewSignInService(authService *auth.Service) *ServiceSignIn {
	return &ServiceSignIn{
		authService: authService,
	}
}

func (p *ServiceSignIn) RenderSignInPage(c *gin.Context) {
	data := gin.H{}
	c.HTML(http.StatusOK, "sign-in", data)
}
func (p *ServiceSignIn) RetrySignIn(c *gin.Context) {
	c.Header("HX-Redirect", "/sign-in")
}
func (p *ServiceSignIn) HandleSignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	tokenPair, err := p.authService.Authenticate(username, password)
	if err != nil {
		logrus.Errorf("Error when verifying user %s", err.Error())
		c.HTML(http.StatusOK, "sign-in_wrong-credentials", gin.H{"Error": "Неверный логин или пароль"})
		return
	}

	c.SetCookie("Bearer", tokenPair.RefreshToken, cookieMaxAge, "/", "", false, true)

	c.HTML(http.StatusOK, "auth-redirect", gin.H{
		"AccessToken": tokenPair.AccessToken,
	})
}
