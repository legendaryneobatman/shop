package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	authService "go-shop/internal/auth/service"
	"go-shop/internal/auth/transport"
	"go-shop/internal/list/entity"
	listService "go-shop/internal/list/service"
	"net/http"
	"strconv"
)

type WebService struct {
	authService    *authService.AuthService
	authMiddleware *transport.AuthMiddleware
	listService    *listService.ListService
}

func NewWebservice(
	authService *authService.AuthService,
	authMiddleware *transport.AuthMiddleware,
	listService *listService.ListService,
) *WebService {
	return &WebService{
		authService:    authService,
		authMiddleware: authMiddleware,
		listService:    listService,
	}
}

func (h *WebService) IndexPage(c *gin.Context) {

	data, err := h.getDataForList(c)

	if err != nil {
		logrus.Errorf("Failed to prepare data for list page, %s", err.Error())
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.HTML(http.StatusOK, "list-elements", data)
		return
	}

	data["Title"] = "Главная"
	c.HTML(http.StatusOK, "main", data)
}

func (h *WebService) getDataForList(c *gin.Context) (gin.H, error) {
	const DefaultOffset = 0
	const Limit = 10

	defaultData := gin.H{
		"lists":      []entity.List{},
		"NextOffset": DefaultOffset + Limit,
		"HasMore":    false,
	}

	data := defaultData

	userId, err := h.authMiddleware.GetUserId(c)
	if err != nil {
		logrus.Errorf("Failed to get user id in index page %s", err.Error())
		return data, err
	}
	offsetStr := c.DefaultQuery("offset", strconv.Itoa(DefaultOffset))
	offset, _ := strconv.Atoi(offsetStr)

	lists, err := h.listService.GetWithPagination(userId, Limit, offset)
	if err != nil {
		logrus.Errorf("Ошибочка: %s", err.Error())
		c.AbortWithStatus(500)
		return data, err
	}
	data["lists"] = lists
	data["NextOffset"] = offset + Limit
	data["HasMore"] = len(lists) == Limit

	return data, nil
}

func (h *WebService) SignInPage(c *gin.Context) {
	data := gin.H{}
	c.HTML(http.StatusOK, "sign-in", data)
}

func (h *WebService) HandleSignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	token, err := h.authService.Authenticate(username, password)
	if err != nil {
		logrus.Errorf("Error when verifying user %s", err.Error())
		c.HTML(http.StatusOK, "content", gin.H{"Error": "Неверный логин или пароль"})
		return
	}

	c.SetCookie("Bearer", token, 3600*24, "/", "", false, true)

	c.Header("HX-Redirect", "/")
}

func (h *WebService) LoadMoreList(c *gin.Context) {
	data, err := h.getDataForList(c)
	if err != nil {
		logrus.Errorf("Failed to load more lists: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "list-elements", data)
}
