package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) OrderHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	orderIndexes, order, err := h.Repository.GetIndexesOrder(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "request_page.html", gin.H{
		"orderIndexes": orderIndexes,
		"order":        order,
		"count":           h.Repository.GetOrderCount(),
	})
}

func (h *Handler) DeleteOrder(ctx *gin.Context){
	idStr := ctx.Param("id")
	orderId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}


	err = h.Repository.DeleteCalculation(orderId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/indexes")
}