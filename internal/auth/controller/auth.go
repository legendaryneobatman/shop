package controller

import (
	"github.com/gin-gonic/gin"
	"go-shop/internal/auth/repository"
	"go-shop/internal/auth/service"
	"go-shop/internal/auth/transport"
	todo "go-shop/internal/user/entity"
	"go-shop/pkg/webtool"
	"net/http"
)

func (a *AuthController) InitRoutes(r *gin.RouterGroup) {
	r.POST("sign-in", a.SignIn)
	r.POST("sign-up", a.SignUp)
	r.POST("refresh", a.Refresh)
	r.POST("logout", a.Logout)
	r.POST("logout-all", a.LogoutAll)
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthController struct {
	RToken *repository.TokenRepository
	SAuth  *service.AuthService
	MAuth  *transport.AuthMiddleware
}

func NewAuthController(
	authService *service.AuthService,
	rt *repository.TokenRepository,
	ma *transport.AuthMiddleware,
) *AuthController {
	return &AuthController{
		SAuth:  authService,
		RToken: rt,
		MAuth:  ma,
	}
}

func (a *AuthController) SignUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := a.SAuth.CreateUser(&input)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
func (a *AuthController) SignIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := a.SAuth.Authenticate(input.Username, input.Password)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	})
}
func (a *AuthController) Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	tokenPair, err := a.SAuth.RefreshTokens(input.RefreshToken)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenPair)

}
func (a *AuthController) Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.SAuth.RevokeToken(input.RefreshToken); err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}
func (a *AuthController) LogoutAll(c *gin.Context) {
	userId, err := a.MAuth.GetUserId(c)

	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, "Not user id found")
		return
	}

	if err := a.SAuth.RevokeAllTokens(userId); err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, "Not user id found")
		return
	}
}
