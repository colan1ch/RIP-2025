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

	indexSearch := ctx.Query("indexSearch")
	if indexSearch == "" {
		indexes, err = h.Repository.GetIndexes()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		indexes, err = h.Repository.GetIndexesByName(indexSearch)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "indexes_page.html", gin.H{
		"time":       time.Now().Format("15:04:05"),
		"indexes":    indexes,
		"indexSearch":      indexSearch,
		"queryId":    h.Repository.GetQueryId(),
		"queryCount": h.Repository.GetQueryCount(h.Repository.GetQueryId()),
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

       ctx.HTML(http.StatusOK, "index_page.html", gin.H{
	       "index": index,
       })
}

func (h *Handler) QueryHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	var indexes []repository.Index = h.Repository.GetQueryIndexes(id)
	query := h.Repository.GetQuery(id)
	
	ctx.HTML(http.StatusOK, "query_page.html", gin.H{
		"queryIndexes": indexes,
		"count":        len(indexes),
		"query": 		query,
		"queryID":      id,
	})
}

func (h *Handler) QueryHandlerUp(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.Atoi(idStr)
    numStr := ctx.Param("num")
    num, _ := strconv.Atoi(numStr)

    query := h.Repository.GetQuery(id)
    if query == nil || num < 1 || num >= len(query.IndexesParametrs) {
        ctx.Redirect(http.StatusSeeOther, "/queries/"+idStr)
        return
    }
    query.IndexesParametrs[num], query.IndexesParametrs[num-1] = query.IndexesParametrs[num-1], query.IndexesParametrs[num]

    ctx.Redirect(http.StatusSeeOther, "/queries/"+idStr)
}

func (h *Handler) QueryHandlerDown(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, _ := strconv.Atoi(idStr)
    numStr := ctx.Param("num")
    num, _ := strconv.Atoi(numStr)

    query := h.Repository.GetQuery(id)
    if query == nil || num < 0 || num >= len(query.IndexesParametrs)-1 {
        ctx.Redirect(http.StatusSeeOther, "/queries/"+idStr)
        return
    }
    query.IndexesParametrs[num], query.IndexesParametrs[num+1] = query.IndexesParametrs[num+1], query.IndexesParametrs[num]

    ctx.Redirect(http.StatusSeeOther, "/queries/"+idStr)
}
