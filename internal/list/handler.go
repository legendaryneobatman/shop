package list

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-shop/internal/auth"
	"go-shop/pkg/webtool"
	"net/http"
	"strconv"
)

type Handler struct {
	Service        *Service
	authService    *auth.Service
	authMiddleware *auth.Middleware
}

func NewHandler(service *Service, authService *auth.Service, authMiddleware *auth.Middleware) *Handler {
	return &Handler{
		Service:        service,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	group := api.Group("/list")

	group.POST("", h.CreateList)
	group.GET("", h.GetLists)
	group.GET(":id", h.GetListByID)
	group.PUT(":id", h.UpdateList)
	group.DELETE(":id", h.DeleteList)
}

func (h *Handler) CreateList(c *gin.Context) {
	userID, err := h.authMiddleware.GetUserID(c)
	if err != nil {
		return
	}

	var input List
	if err := c.BindJSON(&input); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Service.Create(userID, input)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetLists(c *gin.Context) {
	userID, err := h.authMiddleware.GetUserID(c)
	if err != nil {
		return
	}

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr == "" && offsetStr == "" {
		lists, err := h.Service.GetAll(userID)
		if err != nil {
			webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, lists)
		return
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	lists, err := h.Service.GetWithPagination(userID, limit, offset)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *Handler) GetListByID(c *gin.Context) {
	list, err := h.Service.GetByID(c.Param("id"))
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetListWithPagination(c *gin.Context) {
	userID, err := h.authMiddleware.GetUserID(c)
	if err != nil {
		logrus.Errorf("Failed to get user id from context: %s", err.Error())
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		logrus.Errorf("Failed to pars query: limit error: %s", err.Error())
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		logrus.Errorf("Failed to pars query: offset error: %s", err.Error())
	}

	list, err := h.Service.GetWithPagination(userID, limit, offset)
	if err != nil {
		logrus.Errorf("Failed to get list from list service: %s", err.Error())
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) UpdateList(c *gin.Context) {
	updatedList := List{}
	if err := c.BindJSON(&updatedList); err != nil {
		webtool.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	list, err := h.Service.Update(c.Param("id"), updatedList)
	if err != nil {
		webtool.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteList(_ *gin.Context) {

}
