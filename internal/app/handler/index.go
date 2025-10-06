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

	indexSearch := ctx.Query("indexSearch")
	if indexSearch == "" {
		indexes, err = h.Repository.GetIndexes()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	} else {
		indexes, err = h.Repository.GetIndexesByName(indexSearch)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	}
	currentQuery, _ := h.Repository.CheckCurrentQueryDraft(creatorID)

	ctx.HTML(http.StatusOK, "indexes_page.html", gin.H{
		"indexes":       indexes,
		"queryCount": h.Repository.GetQueryCount(),
		"indexSearch":         indexSearch,
		"queryId":    currentQuery.ID,
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

	ctx.HTML(http.StatusOK, "index_page.html", gin.H{
		"index": index,
	})
}


func (h *Handler) AddIndexToQuery(ctx *gin.Context) {
	query, err := h.Repository.GetQueryDraft(h.Repository.GetUser())
	queryId := query.ID
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddIndexToQuery(int(queryId), indexId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/indexes")
}