package handler

import (
	"net/http"
	"strconv"

	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetIndexes(ctx *gin.Context) {
	var indexes []ds.Index
	var err error
	creatorID := h.Repository.GetUser()

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		indexes, err = h.Repository.GetIndexes()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	} else {
		indexes, err = h.Repository.GetIndexesByName(searchQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	}
	currentOrder, _ := h.Repository.CheckCurrentOrderDraft(creatorID)

	ctx.HTML(http.StatusOK, "main_page.html", gin.H{
		"indexes":       indexes,
		"orderCount": h.Repository.GetOrderCount(),
		"query":         searchQuery,
		"orderId":    currentOrder.ID,
	})
}

func (h *Handler) GetIndex(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	index, err := h.Repository.GetIndex(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "product_page.html", gin.H{
		"index": index,
	})
}


func (h *Handler) AddIndexToOrder(ctx *gin.Context) {
	order, err := h.Repository.GetOrderDraft(h.Repository.GetUser())
	orderId := order.ID
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddIndexToOrder(int(orderId), indexId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/indexes")
}