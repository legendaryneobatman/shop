package transport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth/service"
	"go-shop/pkg/webtool"
	"net/http"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (h *AuthMiddleware) UserIdentity(c *gin.Context) {
	token, err := c.Cookie("Bearer")
	if err != nil {
		logrus.Errorf("Failed to extract cookie in userIdentity %s", err.Error())
	}
	if token == "" {
		c.Redirect(http.StatusFound, "/sign-in")
		webtool.NewErrorResponse(c, http.StatusUnauthorized, "empty auth token")
		return
	}

	userId, err := h.authService.ParseToken(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/sign-in")
		webtool.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func (h *AuthMiddleware) GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
