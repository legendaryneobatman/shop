package auth

import (
	"errors"
	"go-shop/internal/models"
	"go-shop/pkg/webtool"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Authenticate(username, password string) (*models.TokenPair, error)
	ParseTokenForUserID(accessToken string) (int, error)
	RefreshTokens(refreshToken string) (*models.TokenPair, error)
	RevokeToken(refreshToken string) error
	RevokeAllTokens(userID int) error
}

type IUserService interface {
	CreateUser(input *models.User) (models.User, error)
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	r.POST("sign-in", webtool.MakeHandler(h.SignIn))
	r.POST("sign-up", webtool.MakeHandler(h.SignUp))
	r.POST("refresh", webtool.MakeHandler(h.Refresh))
	r.POST("logout", webtool.MakeHandler(h.Logout))
	r.POST("logout-all", webtool.MakeHandler(h.LogoutAll))
}

type Handler struct {
	RToken      *RepositoryToken
	service     *Service
	middleware  *Middleware
	errors      *Errors
	userService IUserService
}

func NewHandler(
	service *Service,
	userService IUserService,
	rt *RepositoryToken,
	ma *Middleware,
	errors *Errors,
) *Handler {
	return &Handler{
		service:     service,
		userService: userService,
		RToken:      rt,
		middleware:  ma,
		errors:      errors,
	}
}

func (h *Handler) SignUp(c *gin.Context) *webtool.APIError {
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		return h.errors.UserAlreadyExists
	}

	user, err := h.userService.CreateUser(&models.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	})
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

	c.JSON(http.StatusOK, SignInOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
	return nil
}
func (h *Handler) Refresh(c *gin.Context) *webtool.APIError {
	var input RefreshInput

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

	c.JSON(http.StatusOK, RefreshOutput{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
	return nil
}
func (h *Handler) Logout(c *gin.Context) *webtool.APIError {
	var input LogoutInput

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
	userID, err := h.middleware.GetUserIDCTX(c)

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
