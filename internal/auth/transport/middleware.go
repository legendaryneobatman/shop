package transport

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-shop/internal/auth/service"
	"go-shop/pkg/webtool"
	"net/http"
	"strings"
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

func (h *AuthMiddleware) isAuthenticated(c *gin.Context) bool {
	token := c.GetHeader(authorizationHeader)
	if token == "" {
		webtool.NewErrorResponse(c, http.StatusUnauthorized, "empty auth token")
		return false
	}

	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		webtool.NewErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format")
		return false
	}

	accessToken := parts[1]
	userId, err := h.authService.ParseToken(accessToken)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusUnauthorized, "Invalid or expired access token")
		return false
	}

	c.Set(userCtx, userId)

	return true
}

func (h *AuthMiddleware) WebGuard(c *gin.Context) {
	isAuth := h.isAuthenticated(c)

	if !isAuth {
		println("ZALUPA")
		c.Header("HX-Location", "/sign-in")
		return
	}
}
func (h *AuthMiddleware) ApiGuard(c *gin.Context) {
	isAuth := h.isAuthenticated(c)

	if !isAuth {
		return
	}
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
