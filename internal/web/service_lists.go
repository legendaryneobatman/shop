package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth"
	"go-shop/internal/list"
	"net/http"
	"strconv"
)

type ServiceLists struct {
	authService    *auth.Service
	authMiddleware *auth.Middleware
	listService    *list.Service
}

func NewService(
	authService *auth.Service,
	authMiddleware *auth.Middleware,
	listService *list.Service,
) *ServiceLists {
	return &ServiceLists{
		authService:    authService,
		authMiddleware: authMiddleware,
		listService:    listService,
	}
}

func (s *ServiceLists) RenderListsPage(c *gin.Context) {
	data, err := s.getDataForList(c)

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

func (s *ServiceLists) getDataForList(c *gin.Context) (gin.H, error) {
	const defaultOffset = 0
	const limit = 10

	defaultData := gin.H{
		"lists":      []list.List{},
		"NextOffset": defaultOffset + limit,
		"HasMore":    false,
	}

	data := defaultData

	userID, err := s.authMiddleware.GetUserID(c)
	if err != nil {
		logrus.Errorf("Failed to get user id in index page %s", err.Error())
		return data, err
	}
	offsetStr := c.DefaultQuery("offset", strconv.Itoa(defaultOffset))
	offset, _ := strconv.Atoi(offsetStr)

	lists, err := s.listService.GetWithPagination(userID, limit, offset)
	if err != nil {
		logrus.Errorf("Ошибочка: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return data, err
	}
	data["lists"] = lists
	data["NextOffset"] = offset + limit
	data["HasMore"] = len(lists) == limit

	return data, nil
}

func (s *ServiceLists) LoadMoreList(c *gin.Context) {
	data, err := s.getDataForList(c)
	if err != nil {
		logrus.Errorf("Failed to load more lists: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "list-elements", data)
}
