package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

type Middleware struct {
	service *Service
}

func NewAuthMiddleware(authService *Service) *Middleware {
	return &Middleware{service: authService}
}

func (h *Middleware) isAuthenticated(c *gin.Context) bool {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		logrus.Errorf("Error when get access_token cookie %s", err.Error())
		return false
	}
	if accessToken == "" {
		logrus.Errorf("Empty auth accessToken")
		return false
	}
	userID, err := h.service.ParseTokenForUserID(accessToken)
	println("userID", userID)
	if err != nil || userID == 0 {
		logrus.Errorf("Invalid or expired access accessToken %s", err.Error())
		return false
	}

	c.Set(userCtx, userID)
	println(c.Get(userCtx))

	return true
}

func (h *Middleware) WebGuard(c *gin.Context) {
	isAuth := h.isAuthenticated(c)

	if !isAuth {
		c.Header("HX-Redirect", "/sign-in")
		return
	}
}
func (h *Middleware) APIGuard(c *gin.Context) {
	if !h.isAuthenticated(c) {
		return
	}
}

func (h *Middleware) GetUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		logrus.Errorf("User id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		logrus.Errorf("User id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
