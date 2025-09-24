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
	router.GET("/index/:id", h.GetIndex)
	router.GET("/request/:id", h.OrderHandler)
	// router.GET("/request/:id/up/:num", h.OrderHandlerUp)
	// router.GET("/request/:id/down/:num", h.OrderHandlerDown)
	router.POST("/index/:id/add-to-order", h.AddIndexToOrder)
	router.POST("/request/:id/delete-order", h.DeleteOrder)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("/Users/nachernev/Desktop/rip_lab2/templates/*")
	router.Static("/static", "/Users/nachernev/Desktop/rip_lab2/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
