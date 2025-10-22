package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetIndexes(ctx *gin.Context) {
	var indexes []ds.Index
	var err error

	searchQuery := ctx.Query("index_name")
	if searchQuery == "" {
		indexes, err = h.Repository.GetIndexes()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		indexes, err = h.Repository.GetIndexesByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	resp := make([]apitypes.IndexJSON, 0, len(indexes))
	for _, r := range indexes {
		resp = append(resp, apitypes.IndexToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetIndex(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	index, err := h.Repository.GetIndex(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.IndexToJSON(*index))
}

func (h *Handler) CreateIndex(ctx *gin.Context) {
	var indexJSON apitypes.IndexJSON
	if err := ctx.BindJSON(&indexJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	index, err := h.Repository.CreateIndex(indexJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/indexes/%v", index.ID))
	ctx.JSON(http.StatusCreated, apitypes.IndexToJSON(index))
}

func (h *Handler) DeleteIndex(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteIndex(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"status": "deleted",
	// })
}

func (h *Handler) ChangeIndex(ctx *gin.Context){
	var indexJSON apitypes.IndexJSON
	if err := ctx.BindJSON(&indexJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	index, err := h.Repository.ChangeIndex(id, indexJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.IndexToJSON(index))
}

func (h *Handler) AddIndexToQuery(ctx *gin.Context) {
	query, created, err := h.Repository.GetQueryDraft(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	queryId := query.ID

	indexId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddIndexToQuery(int(queryId), indexId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	
	status := http.StatusOK
	
	if created {
		ctx.Header("Location", fmt.Sprintf("/query/%v", query.ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(query)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, apitypes.QueryToJSON(query, creatorLogin, moderatorLogin))
}

func (h *Handler) UploadImage(ctx *gin.Context) {
	indexId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	index, err := h.Repository.UploadImage(ctx, indexId, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		// "status": "uploaded",
		"index": apitypes.IndexToJSON(index),
	})
}