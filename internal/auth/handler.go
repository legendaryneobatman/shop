package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-shop/pkg/webtool"
	"net/http"
)

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	r.POST("sign-in", MakeHandler(h.SignIn))
	r.POST("sign-up", MakeHandler(h.SignUp))
	r.POST("refresh", MakeHandler(h.Refresh))
	r.POST("logout", MakeHandler(h.Logout))
	r.POST("logout-all", MakeHandler(h.LogoutAll))
}

func MakeHandler(handlerWithError func(c *gin.Context) *webtool.APIError) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handlerWithError(c)
		if err.Error != nil {
			webtool.NewErrorResponse(c, err.HTTPStatus, err.Error.Error())
			return
		}
	}
}

type Handler struct {
	RToken     *RepositoryToken
	service    *Service
	middleware *Middleware
	errors     *Errors
}

func NewAuthController(
	authService *Service,
	rt *RepositoryToken,
	ma *Middleware,
	errors *Errors,
) *Handler {
	return &Handler{
		service:    authService,
		RToken:     rt,
		middleware: ma,
		errors:     errors,
	}
}

func (h *Handler) SignUp(c *gin.Context) *webtool.APIError {
	var input SignUpInput

	if err := c.BindJSON(&input); err != nil {
		return h.errors.UserAlreadyExists
	}

	user, err := h.service.CreateUser(&input)
	if err != nil {
		if errors.Is(err, h.errors.UserAlreadyExists.Error) {
			return h.errors.UserAlreadyExists
		}
	}

	c.JSON(http.StatusOK, SignUpOutput{
		ID: user.ID,
	})
	return nil
}
func (h *Handler) SignIn(c *gin.Context) *webtool.APIError {
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		return h.errors.InvalidCredentials
	}

	token, err := h.service.Authenticate(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, h.errors.InvalidCredentials.Error) {
			return h.errors.InvalidCredentials
		}
	}

	c.JSON(http.StatusOK, map[string]any{
		"accessToken":  token.AccessToken,
		"refreshToken": token.RefreshToken,
	})
	return nil
}
func (h *Handler) Refresh(c *gin.Context) *webtool.APIError {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		return h.errors.InvalidCredentials
	}
	tokenPair, err := h.service.RefreshTokens(input.RefreshToken)
	if err != nil {
		if errors.Is(err, h.errors.NoTokenFound.Error) {
			return h.errors.NoTokenFound
		}
		return h.errors.InvalidCredentials
	}

	c.JSON(http.StatusOK, tokenPair)
	return nil
}
func (h *Handler) Logout(c *gin.Context) *webtool.APIError {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		return h.errors.InvalidCredentials
	}

	if err := h.service.RevokeToken(input.RefreshToken); err != nil {
		if errors.Is(err, h.errors.NoTokenFound.Error) {
			return h.errors.NoTokenFound
		}
		return h.errors.InvalidCredentials
	}
	return nil
}
func (h *Handler) LogoutAll(c *gin.Context) *webtool.APIError {
	userID, err := h.middleware.GetUserID(c)

	if err != nil {
		return h.errors.UserUnauthorized
	}

	if err := h.service.RevokeAllTokens(userID); err != nil {
		if errors.Is(err, h.errors.NoTokenFound.Error) {
			return h.errors.NoTokenFound
		}
		return h.errors.InvalidCredentials
	}
	return nil
}
