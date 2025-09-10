package handler

import (
	"LAB1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetIndexes(ctx *gin.Context) {
	var indexes []repository.Index
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		indexes, err = h.Repository.GetIndexes()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		indexes, err = h.Repository.GetIndexesByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "main_page.html", gin.H{
		"time":          time.Now().Format("15:04:05"),
		"indexes":       indexes,
		"query":         searchQuery,
		"researchId":    h.Repository.GetResearchId(),
		"researchCount": h.Repository.GetResearchCount(h.Repository.GetResearchId()),
	})
}

func (h *Handler) GetIndex(ctx *gin.Context) {
       idStr := ctx.Param("id")
       id, err := strconv.Atoi(idStr)
       if err != nil {
	       logrus.Error(err)
       }

       index, err := h.Repository.GetIndex(id)
       if err != nil {
	       logrus.Error(err)
       }

       ctx.HTML(http.StatusOK, "product_page.html", gin.H{
	       "index": index,
       })
}

func (h *Handler) RequestHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	var indexes []repository.Index
	indexes = h.Repository.GetResearchIndexes(id)

	ctx.HTML(http.StatusOK, "request_page.html", gin.H{
		"researchIndexes": indexes,
		"count":           len(indexes),
	})
}
