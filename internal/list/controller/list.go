package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	auth "go-shop/internal/auth/service"
	"go-shop/internal/auth/transport"
	"go-shop/internal/list/entity"
	"go-shop/pkg/webtool"
	"net/http"
	"strconv"
)

type ListService interface {
	Create(userId int, list entity.List) (int, error)
	GetAll(userId int) ([]entity.List, error)
	GetWithPagination(userId int, limit int, offset int) ([]entity.List, error)
	GetById(listId string) (entity.List, error)
	Update(listId string, input entity.List) (entity.List, error)
}

type ListController struct {
	listService    ListService
	authService    *auth.AuthService
	authMiddleware *transport.AuthMiddleware
}

func NewListController(listService ListService, authService *auth.AuthService, authMiddleware *transport.AuthMiddleware) *ListController {
	return &ListController{
		listService:    listService,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (listController *ListController) InitRoutes(api *gin.RouterGroup) {
	auth := api.Group("/list")

	auth.POST("", listController.CreateList)
	auth.GET("", listController.GetLists)
	auth.GET(":id", listController.GetListById)
	auth.PUT(":id", listController.UpdateList)
	auth.DELETE(":id", listController.DeleteList)
}

func (listController *ListController) CreateList(c *gin.Context) {
	userId, err := listController.authMiddleware.GetUserId(c) // Не смотреть сильно
	if err != nil {
		return
	}

	var input entity.List
	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := listController.listService.Create(userId, input)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (listController *ListController) GetLists(c *gin.Context) {
	userId, err := listController.authMiddleware.GetUserId(c)
	if err != nil {
		return
	}

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr == "" && offsetStr == "" {
		lists, err := listController.listService.GetAll(userId)
		if err != nil {
			webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, lists)
		return
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	lists, err := listController.listService.GetWithPagination(userId, limit, offset)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (listController *ListController) GetListById(c *gin.Context) {
	list, err := listController.listService.GetById(c.Param("id"))
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (listController *ListController) GetListWithPagination(c *gin.Context) {
	userId, err := listController.authMiddleware.GetUserId(c)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		logrus.Errorf("Failed to pars query: limit error: %s", err.Error())
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		logrus.Errorf("Failed to pars query: offset error: %s", err.Error())
	}
	if err != nil {
		logrus.Errorf("No user id provided for GetListsWithPagination error: %s", err.Error())
	}

	list, err := listController.listService.GetWithPagination(userId, limit, offset)
	if err != nil {
		logrus.Errorf("Failed to get list from list service: %s", err.Error())
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, list)
	return
}

func (listController *ListController) UpdateList(c *gin.Context) {
	updatedList := entity.List{}
	if err := c.BindJSON(&updatedList); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if list, err := listController.listService.Update(c.Param("id"), updatedList); err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	} else {
		c.JSON(http.StatusOK, list)
		return
	}
}

func (listController *ListController) DeleteList(c *gin.Context) {

}
