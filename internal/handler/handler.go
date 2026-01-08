package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	authController "go-shop/internal/auth/controller"
	authRepository "go-shop/internal/auth/repository"
	authService "go-shop/internal/auth/service"
	"go-shop/internal/auth/transport"
	listController "go-shop/internal/list/controller"
	listRepository "go-shop/internal/list/repository"
	listService "go-shop/internal/list/service"
	webController "go-shop/internal/web/controller"
	webService "go-shop/internal/web/service"
)

type Handler struct {
	db             *sqlx.DB
	listService    *listService.ListService
	authMiddleware *transport.AuthMiddleware
	authService    *authService.AuthService
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Init(router *gin.Engine) {
	//auth
	authRepo := authRepository.NewAuthRepository(h.db)
	authServ := authService.NewAuthService(authRepo)
	authMiddleware := transport.NewAuthMiddleware(authServ)
	authCtrl := authController.NewAuthController(authServ)

	h.authMiddleware = authMiddleware
	h.authService = authServ

	webGroup := router.Group("", h.authMiddleware.UserIdentity)
	apiPGroup := router.Group("/api", h.authMiddleware.UserIdentity)
	apiGroup := router.Group("/api")
	authGroup := router.Group("/")

	//list
	listRepo := listRepository.NewListRepository(h.db)
	listServ := listService.NewListService(listRepo)
	listCtrl := listController.NewListController(listServ, authServ, authMiddleware)
	h.listService = listServ

	//todo

	//статика
	webServ := webService.NewWebservice(authServ, authMiddleware, listServ)
	webCtrl := webController.NewWebController(webServ, authMiddleware, listServ, authServ)

	authCtrl.InitRoutes(apiGroup)
	listCtrl.InitRoutes(apiPGroup)
	webCtrl.InitRoutes(webGroup, authGroup)
}
