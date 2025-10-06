package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) QueryHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	queryIndexes, query, err := h.Repository.GetIndexesQuery(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.HTML(http.StatusOK, "query_page.html", gin.H{
		"queryIndexes": queryIndexes,
		"query":        query,
		"count":           h.Repository.GetQueryCount(),
	})
}

func (h *Handler) DeleteQuery(ctx *gin.Context){
	idStr := ctx.Param("id")
	queryId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}


	err = h.Repository.DeleteQuery(queryId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/indexes")
}