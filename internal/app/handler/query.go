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
		// h.errorHandler(ctx, http.StatusInternalServerError, err)
		ctx.Redirect(http.StatusFound, "/indexes/")
		return
	}
	ctx.HTML(http.StatusOK, "query_page.html", gin.H{
		"queryIndexes": queryIndexes,
		"query":        query,
		"count":           h.Repository.GetQueryCount(),
		"count1": h.Repository.GetQueryCount() - 1,
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

func (h *Handler) UpdateRowsCount(ctx *gin.Context){
	idStr := ctx.Param("id")
	indexIdStr := ctx.Param("indexId")
	rowsCountStr := ctx.Query("rowsCount")
	queryId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	indexId, err := strconv.Atoi(indexIdStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	rowsCount, err := strconv.Atoi(rowsCountStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.UpdateRowsCount(queryId, indexId, rowsCount)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/query/"+idStr)
}

func (h *Handler) UpdateRecievedRows(ctx *gin.Context){
	idStr := ctx.Param("id")
	indexIdStr := ctx.Param("indexId")
	recievedRowsStr := ctx.Query("recievedRows")
	queryId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	indexId, err := strconv.Atoi(indexIdStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	recievedRows, err := strconv.Atoi(recievedRowsStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.UpdateRecievedRows(queryId, indexId, recievedRows)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/query/"+idStr)
}

func (h *Handler) UpdateTableField(ctx *gin.Context){
	idStr := ctx.Param("id")
	indexIdStr := ctx.Param("indexId")
	tableField := ctx.Query("tableField")
	queryId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	indexId, err := strconv.Atoi(indexIdStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.UpdateTableField(queryId, indexId, tableField)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/query/"+idStr)
}
