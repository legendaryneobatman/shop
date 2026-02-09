package bootstrap

import (
	"go-shop/internal/auth"
	"go-shop/internal/category"
	"go-shop/internal/docs"
	"go-shop/internal/list"
	"go-shop/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Init(router *gin.Engine) {
	tokenRepo := auth.NewRepositoryToken(h.db)
	userRepo := user.NewRepository(h.db)
	listRepo := list.NewRepository(h.db)
	categoryRepo := category.NewRepository(h.db)

	userServ := user.NewService(userRepo)
	authServ := auth.NewAuthService(userRepo, tokenRepo)
	listServ := list.NewListService(listRepo)
	categoryServ := category.NewService(categoryRepo)

	authMiddleware := auth.NewAuthMiddleware(authServ)

	apiProtectedGroup := router.Group("/api", authMiddleware.APIGuard)
	apiPublicGroup := router.Group("/api")
	authGroup := router.Group("/")

	docsCtrl := docs.NewHandler()
	authCtrl := auth.NewHandler(authServ, userServ, tokenRepo, authMiddleware)
	listCtrl := list.NewHandler(listServ, authMiddleware)
	categoryCtrl := category.NewHandler(categoryServ)

	authCtrl.InitRoutes(apiPublicGroup)
	listCtrl.InitRoutes(apiProtectedGroup)
	docsCtrl.InitRoutes(authGroup)
	categoryCtrl.InitRoutes(apiProtectedGroup)
}
