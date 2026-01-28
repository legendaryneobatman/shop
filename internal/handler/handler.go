package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go-shop/internal/auth"
	"go-shop/internal/list"
	"go-shop/internal/web"
)

type Handler struct {
	db             *sqlx.DB
	listService    *list.Service
	authMiddleware *auth.Middleware
	authService    *auth.Service
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Init(router *gin.Engine) {
	tokenRepo := auth.NewRepositoryToken(h.db)

	// auth
	authErrors := auth.NewAuthErrors()
	authRepo := auth.NewRepositoryAuth(h.db)
	authServ := auth.NewAuthService(authRepo, tokenRepo, authErrors)
	authMiddleware := auth.NewAuthMiddleware(authServ)
	authCtrl := auth.NewAuthController(authServ, tokenRepo, authMiddleware, authErrors)

	h.authMiddleware = authMiddleware
	h.authService = authServ

	webGroup := router.Group("", h.authMiddleware.WebGuard)
	apiPGroup := router.Group("/api", h.authMiddleware.APIGuard)
	apiGroup := router.Group("/api")
	authGroup := router.Group("/")

	// list
	listRepo := list.NewRepository(h.db)
	listServ := list.NewListService(listRepo)
	listCtrl := list.NewHandler(listServ, authServ, authMiddleware)
	h.listService = listServ

	// todo

	// статика
	webServ := web.NewService(authServ, authMiddleware, listServ)
	signInServ := web.NewSignInService(authServ)
	webCtrl := web.NewHandler(webServ, signInServ)

	authCtrl.InitRoutes(apiGroup)
	listCtrl.InitRoutes(apiPGroup)
	webCtrl.InitRoutes(webGroup, authGroup)
}
