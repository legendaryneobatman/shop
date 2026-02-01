package auth

import (
	"errors"
	"net/http"
	"strings"

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
	accessToken := h.getAccessTokenFromHeader(c)
	if accessToken == "" {
		logrus.Errorf("Empty auth accessToken")
		return false
	}
	userID, err := h.service.ParseTokenForUserID(accessToken)
	if err != nil || userID == 0 {
		logrus.Errorf("Invalid or expired access accessToken %s", err.Error())
		return false
	}

	c.Set(userCtx, userID)

	return true
}

func (h *Middleware) getAccessTokenFromCookies(c *gin.Context) string {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		logrus.Errorf("Error when get access_token cookie %s", err.Error())
		return ""
	}
	return accessToken
}
func (h *Middleware) getAccessTokenFromHeader(c *gin.Context) string {
	header := c.GetHeader(authorizationHeader)
	logrus.Debugln("header", header)
	accessToken := strings.TrimPrefix(header, "Bearer ")
	if accessToken == "" {
		logrus.Errorf("Empty auth accessToken")
		return ""
	}
	return accessToken
}
func (h *Middleware) APIGuard(c *gin.Context) {
	if !h.isAuthenticated(c) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (h *Middleware) GetUserIDCTX(c *gin.Context) (int, error) {
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
