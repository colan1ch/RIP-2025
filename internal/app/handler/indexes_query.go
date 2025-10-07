package handler

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteIndexFromQuery(ctx *gin.Context) {
	queryId, err := strconv.Atoi(ctx.Param("query_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexId, err := strconv.Atoi(ctx.Param("index_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	query, err := h.Repository.DeleteIndexFromQuery(queryId, indexId)
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

func (h *Handler) ChangeIndexQuery(ctx *gin.Context) {
	queryId, err := strconv.Atoi(ctx.Param("query_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexId, err := strconv.Atoi(ctx.Param("index_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var indexQueryJSON apitypes.IndexesQueryJSON
	if err := ctx.BindJSON(&indexQueryJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	indexQuery, err := h.Repository.ChangeIndexQuery(queryId, indexId, indexQueryJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.IndexesQueryToJSON(indexQuery))
}