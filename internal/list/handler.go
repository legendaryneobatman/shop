package list

import (
	"go-shop/internal/models"
	"go-shop/pkg/webtool"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IAuthMiddleware interface {
	GetUserIDCTX(c *gin.Context) (int, error)
}

type Handler struct {
	Service        *Service
	authMiddleware IAuthMiddleware
}

func NewHandler(service *Service, authMiddleware IAuthMiddleware) *Handler {
	return &Handler{
		Service:        service,
		authMiddleware: authMiddleware,
	}
}

func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	group := api.Group("/list")

	group.POST("", webtool.MakeHandler(h.CreateList))
	group.GET("", webtool.MakeHandler(h.GetLists))
	group.GET(":id", webtool.MakeHandler(h.GetListByID))
	group.PUT(":id", webtool.MakeHandler(h.UpdateList))
	group.DELETE(":id", webtool.MakeHandler(h.DeleteList))
}

func (h *Handler) CreateList(c *gin.Context) *webtool.APIError {
	userID, err := h.authMiddleware.GetUserIDCTX(c)
	if err != nil {
		return NoUserFound(err)
	}

	var input models.List
	if err := c.BindJSON(&input); err != nil {
		return BadRequest(err)
	}

	id, err := h.Service.Create(userID, input)
	if err != nil {
		return CantCreate(err)
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})

	return nil
}

func (h *Handler) GetLists(c *gin.Context) *webtool.APIError {
	userID, err := h.authMiddleware.GetUserIDCTX(c)
	if err != nil {
		return NoUserFound(err)
	}

	lists, err := h.Service.GetAll(userID)
	if err != nil {
		return CantFindElements(err)
	}

	c.JSON(http.StatusOK, lists)
	return nil
}

func (h *Handler) GetListWithPagination(c *gin.Context) *webtool.APIError {
	userID, err := h.authMiddleware.GetUserIDCTX(c)
	if err != nil {
		return NoUserFound(err)
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return InvalidPagination(err)
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return InvalidPagination(err)
	}

	list, err := h.Service.GetWithPagination(userID, limit, offset)
	if err != nil {
		return CantFindElements(err)
	}

	c.JSON(http.StatusOK, list)

	return nil
}

func (h *Handler) GetListByID(c *gin.Context) *webtool.APIError {
	list, err := h.Service.GetByID(c.Param("id"))
	if err != nil {
		return CantFindElements(err)
	}

	c.JSON(http.StatusOK, list)

	return nil
}

func (h *Handler) UpdateList(c *gin.Context) *webtool.APIError {
	updatedList := models.List{}
	if err := c.BindJSON(&updatedList); err != nil {
		return BadRequest(err)
	}
	list, err := h.Service.Update(c.Param("id"), updatedList)
	if err != nil {
		return CantFindElements(err)
	}

	c.JSON(http.StatusOK, list)

	return nil
}

func (h *Handler) DeleteList(c *gin.Context) *webtool.APIError {
	params := []string{c.Param("id")}
	ids, err := h.Service.Delete(params)
	if err != nil {
		return CantFindElements(err)
	}

	c.JSON(http.StatusOK, ids)

	return nil
}
