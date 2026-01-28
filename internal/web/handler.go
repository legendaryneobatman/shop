package web

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	_webService   *ServiceLists
	signInService *ServiceSignIn
}

func NewHandler(
	webService *ServiceLists,
	signInService *ServiceSignIn,
) *Handler {
	return &Handler{
		_webService:   webService,
		signInService: signInService,
	}
}

func (w *Handler) InitRoutes(pg *gin.RouterGroup, upg *gin.RouterGroup) {
	// тут разные должны быть группы, для ауфа непротектед, для листов протектед
	upg.GET("/sign-in", w.signInService.RenderSignInPage)
	upg.POST("/login", w.signInService.HandleSignIn)
	upg.POST("/sign-in-retry", w.signInService.RetrySignIn)

	pg.GET("/", w._webService.RenderListsPage)
	pg.GET("/load-list", w._webService.LoadMoreList)
}
