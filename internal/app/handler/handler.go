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
		"time":       time.Now().Format("15:04:05"),
		"indexes":    indexes,
		"query":      searchQuery,
		"OrderId":    h.Repository.GetOrderId(),
		"orderCount": h.Repository.GetOrderCount(h.Repository.GetOrderId()),
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

func (h *Handler) OrderHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	var indexes []repository.Index = h.Repository.GetOrderIndexes(id)
	order := h.Repository.GetOrder(id)
	
	ctx.HTML(http.StatusOK, "request_page.html", gin.H{
		"orderIndexes": indexes,
		"count":        len(indexes),
		"order": 		order,
		"orderID":      id,
	})
}

func (h *Handler) OrderHandlerUp(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.Atoi(idStr)
    numStr := ctx.Param("num")
    num, _ := strconv.Atoi(numStr)

    order := h.Repository.GetOrder(id)
    if order == nil || num < 1 || num >= len(order.IndexesParametrs) {
        ctx.Redirect(http.StatusSeeOther, "/request/"+idStr)
        return
    }
    order.IndexesParametrs[num], order.IndexesParametrs[num-1] = order.IndexesParametrs[num-1], order.IndexesParametrs[num]

    ctx.Redirect(http.StatusSeeOther, "/request/"+idStr)
}

func (h *Handler) OrderHandlerDown(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.Atoi(idStr)
    numStr := ctx.Param("num")
    num, _ := strconv.Atoi(numStr)

    order := h.Repository.GetOrder(id)
    if order == nil || num < 0 || num >= len(order.IndexesParametrs)-1 {
        ctx.Redirect(http.StatusSeeOther, "/request/"+idStr)
        return
    }
    order.IndexesParametrs[num], order.IndexesParametrs[num+1] = order.IndexesParametrs[num+1], order.IndexesParametrs[num]

    ctx.Redirect(http.StatusSeeOther, "/request/"+idStr)
}