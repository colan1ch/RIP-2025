package handler

import (
	"LAB1/internal/app/repository"
	// "net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	router.GET("api/indexes", h.GetIndexes)
	router.GET("api/indexes/:id", h.GetIndex)
	router.POST("/api/indexes/:id/add-to-query", h.AddIndexToQuery)
	router.POST("/api/indexes/create-index", h.CreateIndex)
	router.DELETE("/api/indexes/:id/delete-index", h.DeleteIndex)
	router.PUT("/api/indexes/:id/change-index", h.ChangeIndex)
	router.POST("/api/indexes/:id/create-image", h.UploadImage)

	// router.GET("/queries/:id", h.QueryHandler)
	// router.GET("/queries/:id/up/:num", h.QueryHandlerUp)
	// router.GET("/queries/:id/down/:num", h.QueryHandlerDown)
	// router.GET("/queries/:id/update-rows-count/:indexId", h.UpdateRowsCount)
	// router.GET("/queries/:id/update-cardinality/:indexId", h.UpdateCardinality)
	// router.GET("/queries/:id/update-table-field/:indexId", h.UpdateTableField)

	// router.POST("api/queries/:id/delete-query", h.DeleteQuery)
	router.GET("/api/queries/query-cart", h.GetQueryCart)	
	router.GET("/api/queries", h.GetQueries)
	router.GET("/api/queries/:id", h.GetQuery)
	router.PUT("/api/queries/:id/change-query", h.ChangeQuery)
	router.PUT("/api/queries/:id/form", h.FormQuery)
	router.PUT("/api/queries/:id/finish", h.ModerateQuery)
	router.DELETE("/api/queries/:id/delete-query", h.DeleteQuery)

	router.DELETE("/api/indexes_query/:index_id/:query_id", h.DeleteIndexFromQuery)
	router.PUT("/api/indexes_query/:index_id/:query_id", h.ChangeIndexQuery)

	router.POST("/api/users/sign-up", h.CreateUser)
	router.GET("/api/users/profile", h.GetProfile)
	router.PUT("/api/users/profile", h.ChangeProfile)
	router.POST("/api/users/sign-in", h.SignIn)
	router.POST("/api/users/sign-out", h.SignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("/Users/nachernev/Desktop/rip_lab3/templates/*")
	router.Static("/static", "/Users/nachernev/Desktop/rip_lab3/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
