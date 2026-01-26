package controller

import (
	"github.com/gin-gonic/gin"
	"go-shop/internal/web/service"
)

type WebController struct {
	_webService   *service.WebService
	signInService *service.SignInService
}

func NewWebController(
	webService *service.WebService,
	signInService *service.SignInService,
) *WebController {
	return &WebController{
		_webService:   webService,
		signInService: signInService,
	}
}

func (w *WebController) InitRoutes(pg *gin.RouterGroup, upg *gin.RouterGroup) {
	// тут разные должны быть группы, для ауфа непротектед, для листов протектед
	upg.GET("/sign-in", w.signInService.RenderSignInPage)
	upg.POST("/login", w.signInService.HandleSignIn)
	upg.POST("/sign-in-retry", w.signInService.RetrySignIn)

	pg.GET("/", w._webService.RenderListsPage)
	pg.GET("/load-list", w._webService.LoadMoreList)
}
