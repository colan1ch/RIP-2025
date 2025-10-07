package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"LAB1/internal/app/api_types"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetQueries(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}
	fmt.Println(fromDate)

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	queries, err := h.Repository.GetQueries(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]apitypes.QueryJSON, 0, len(queries))
	for _, c := range queries {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, apitypes.QueryToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetQueryCart(ctx *gin.Context){
	indexesCount := h.Repository.GetQueryCount(h.Repository.GetUserID())

	if indexesCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"indexes_count": indexesCount,
		})
		return
	}

	query, err := h.Repository.CheckCurrentQueryDraft(h.Repository.GetUserID())
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":          "no_draft",
				"indexes_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          query.ID,
		"indexes_count": h.Repository.GetQueryCount(query.CreatorID),
	})
}

func (h *Handler) GetQuery(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexes, query, err := h.Repository.GetQueryIndexes(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	resp := make([]apitypes.IndexJSON, 0, len(indexes))
	for _, r := range indexes {
		resp = append(resp, apitypes.IndexToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(query)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	indexesQuery, _ := h.Repository.GetIndexesQueries(query.ID)
	
	resp2 := make([]apitypes.IndexesQueryJSON, 0, len(indexesQuery))
	for _, r := range indexesQuery {
		resp2 = append(resp2, apitypes.IndexesQueryToJSON(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"query": apitypes.QueryToJSON(query, creatorLogin, moderatorLogin),
		"indexes":   resp,
		"indexesQuery": resp2,
	})
}

func (h *Handler) FormQuery(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "formed"

	query, err := h.Repository.FormQuery(id, status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(query)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.QueryToJSON(query, creatorLogin, moderatorLogin))
}

func (h *Handler) ChangeQuery(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var queryJSON apitypes.QueryJSON
	if err := ctx.BindJSON(&queryJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	query, err := h.Repository.ChangeQuery(id, queryJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(query)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.QueryToJSON(query, creatorLogin, moderatorLogin))
}

func (h *Handler) DeleteQuery(ctx *gin.Context){
	idStr := ctx.Param("id")
	queryId, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "deleted"
	
	_, err = h.Repository.FormQuery(queryId, status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Query deleted"})
}

func (h *Handler) ModerateQuery(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON apitypes.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	query, err := h.Repository.ModerateQuery(id, statusJSON.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(query)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.QueryToJSON(query, creatorLogin, moderatorLogin))
}