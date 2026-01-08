package controller

import (
	"github.com/gin-gonic/gin"
	authService "go-shop/internal/auth/service"
	"go-shop/internal/auth/transport"
	listService "go-shop/internal/list/service"
	"go-shop/internal/web/service"
)

type WebController struct {
	_webService *service.WebService

	authMiddleware *transport.AuthMiddleware
	listService    *listService.ListService
	authService    *authService.AuthService
}

func NewWebController(
	_webService *service.WebService,

	authMiddleware *transport.AuthMiddleware,
	listService *listService.ListService,
	authService *authService.AuthService,
) *WebController {
	return &WebController{
		_webService: _webService,

		authMiddleware: authMiddleware,
		listService:    listService,
		authService:    authService,
	}
}

func (w *WebController) InitRoutes(pg *gin.RouterGroup, upg *gin.RouterGroup) {
	// тут разные должны быть группы, для ауфа непротектед, для листов протектед
	upg.GET("/sign-in", w.signInPage)
	upg.POST("/login", w.signInHelper)

	pg.GET("/", w.listPage)
	pg.GET("/list", w.listPageHelper)
}

func (w *WebController) signInPage(c *gin.Context) {
	w._webService.SignInPage(c)
}
func (w *WebController) signInHelper(c *gin.Context) {
	w._webService.Login(c)
}

func (w *WebController) listPage(c *gin.Context) {
	w._webService.IndexPage(c)
}
func (w WebController) listPageHelper(c *gin.Context) {
	w._webService.LoadMoreList(c)
}
