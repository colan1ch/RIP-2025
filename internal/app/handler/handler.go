package handler

import (
	"LAB1/internal/app/repository"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.Use(CORSMiddleware())
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/sign-up", h.SignUp)
	unauthorized.GET("/indexes", h.GetIndexes)
	unauthorized.GET("/indexes/:id", h.GetIndex)
	unauthorized.POST("/users/sign-in", h.SignIn)

	unauthorized.POST("/queries/:id/update-query-result", h.UpdateQueryResult)

	optionalauthorized := api.Group("/")
	optionalauthorized.Use(h.WithOptionalAuthCheck())
	optionalauthorized.GET("/queries/query-cart", h.GetQueryCart)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))

	authorized.POST("/indexes/create-index", h.CreateIndex)
	authorized.DELETE("/indexes/:id/delete-index", h.DeleteIndex)
	authorized.PUT("/indexes/:id/change-index", h.ChangeIndex)
	authorized.POST("/indexes/:id/add-to-query", h.AddIndexToQuery)
	authorized.POST("/indexes/:id/create-image", h.UploadImage)

	authorized.GET("/queries", h.GetQueries)
	authorized.GET("/queries/:id", h.GetQuery)
	authorized.PUT("/queries/:id/change-query", h.ChangeQuery)
	authorized.PUT("/queries/:id/form", h.FormQuery)
	authorized.PUT("/queries/:id/finish", h.ModerateQuery)
	authorized.DELETE("/queries/:id/delete-query", h.DeleteQuery)

	authorized.DELETE("/indexes_query/:index_id/:query_id", h.DeleteIndexFromQuery)
	authorized.PUT("/indexes_query/:index_id/:query_id", h.ChangeIndexQuery)

	authorized.GET("/users/:login/profile", h.GetProfile)
	authorized.PUT("/users/:login/profile", h.ChangeProfile)
	authorized.POST("/users/sign-out", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/query-time-calc/:id/moderate", h.ModerateQuery)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "/Users/nachernev/Desktop/rip_lab4/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())

	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		// "status":      "error",
		"description": errorMessage,
	})
}
