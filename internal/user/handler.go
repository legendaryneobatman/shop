package user

import (
	"go-shop/pkg/webtool"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	group := r.Group("/user")

	group.GET("", webtool.MakeHandler(h.Get))
	group.GET(":id", webtool.MakeHandler(h.GetById))
	group.POST(":id", webtool.MakeHandler(h.Create))
	group.PATCH(":id", webtool.MakeHandler(h.Edit))
	group.DELETE(":id", webtool.MakeHandler(h.Delete))
}

func (h *Handler) Get(c *gin.Context) *webtool.APIError     {



	return nil
}
func (h *Handler) GetCurrent(c *gin.Context) *webtool.APIError { return nil }
func (h *Handler) GetById(c *gin.Context) *webtool.APIError { return nil }
func (h *Handler) Create(c *gin.Context) *webtool.APIError  { return nil }
func (h *Handler) Edit(c *gin.Context) *webtool.APIError    { return nil }
func (h *Handler) Delete(c *gin.Context) *webtool.APIError  { return nil }
