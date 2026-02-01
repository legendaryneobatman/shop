package docs

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	r.GET("/docs", h.ServeDocsYAML)
}

func (h *Handler) ServeDocsYAML(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Methods", "GET")
	context.File("swagger.yaml")
}
