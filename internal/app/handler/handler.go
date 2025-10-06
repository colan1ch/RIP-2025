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
	router.GET("/indexes", h.GetIndexes)
	router.GET("/indexes/:id", h.GetIndex)
	router.GET("/query/:id", h.QueryHandler)
	// router.GET("/request/:id/up/:num", h.QueryHandlerUp)
	// router.GET("/request/:id/down/:num", h.QueryHandlerDown)
	router.POST("/indexes/:id/add-to-query", h.AddIndexToQuery)
	router.POST("/query/:id/delete-query", h.DeleteQuery)
	router.GET("/query/:id/update-rows-count/:indexId", h.UpdateRowsCount)
	router.GET("/query/:id/update-cardinality/:indexId", h.UpdateCardinality)
	router.GET("/query/:id/update-table-field/:indexId", h.UpdateTableField)

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
