package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth/service"
	todo "go-shop/internal/user/entity"
	"go-shop/pkg/webtool"
	"net/http"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) InitRoutes(r *gin.RouterGroup) {
	r.POST("sign-in", controller.SignIn)
	r.POST("sign-up", controller.SignUp)
}

func (controller *AuthController) SignUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := controller.authService.CreateUser(input)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (controller *AuthController) SignIn(c *gin.Context) {
	var input signInInput

	logrus.Debugf(input.Username, input.Password)
	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := controller.authService.VerifyUser(input.Username, input.Password)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
