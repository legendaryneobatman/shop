package bootstrap

import (
	"go-shop/internal/auth"
	"go-shop/internal/docs"
	"go-shop/internal/list"
	"go-shop/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	// docs
	docsCtrl := docs.NewHandler()

	// user
	userRepo := user.NewRepository(h.db)
	userServ := user.NewService(userRepo)

	// auth
	authErrors := auth.NewAuthErrors()
	authRepo := user.NewRepository(h.db)
	authServ := auth.NewAuthService(authRepo, tokenRepo, authErrors)
	authMiddleware := auth.NewAuthMiddleware(authServ)
	authCtrl := auth.NewHandler(authServ, userServ, tokenRepo, authMiddleware, authErrors)

	h.authMiddleware = authMiddleware
	h.authService = authServ

	apiPGroup := router.Group("/api", h.authMiddleware.APIGuard)
	apiGroup := router.Group("/api")
	authGroup := router.Group("/")

	// list
	listErrors := list.NewErrors()
	listRepo := list.NewRepository(h.db)
	listServ := list.NewListService(listRepo)
	listCtrl := list.NewHandler(listServ, listErrors, authMiddleware)
	h.listService = listServ

	authCtrl.InitRoutes(apiGroup)
	listCtrl.InitRoutes(apiPGroup)
	docsCtrl.InitRoutes(authGroup)
}
