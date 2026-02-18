package bootstrap

import (
	"shop/internal/auth"
	"shop/internal/category"
	"shop/internal/docs"
	"shop/internal/list"
	"shop/internal/product"
	"shop/internal/user"

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
	productRepo := product.NewRepository(h.db, categoryRepo)

	userServ := user.NewService(userRepo)
	authServ := auth.NewAuthService(userRepo, tokenRepo)
	listServ := list.NewListService(listRepo)
	categoryServ := category.NewService(categoryRepo)
	productServ := product.NewService(productRepo)

	authMiddleware := auth.NewAuthMiddleware(authServ)

	apiProtectedGroup := router.Group("/api", authMiddleware.APIGuard)
	apiPublicGroup := router.Group("/api")
	authGroup := router.Group("/")

	docsCtrl := docs.NewHandler()
	authCtrl := auth.NewHandler(authServ, userServ, tokenRepo, authMiddleware)
	listCtrl := list.NewHandler(listServ, authMiddleware)
	categoryCtrl := category.NewHandler(categoryServ)
	productCtrl := product.NewHandler(productServ)

	authCtrl.InitRoutes(apiPublicGroup)
	listCtrl.InitRoutes(apiProtectedGroup)
	docsCtrl.InitRoutes(authGroup)
	categoryCtrl.InitRoutes(apiProtectedGroup)
	productCtrl.InitRoutes(apiProtectedGroup)
}
