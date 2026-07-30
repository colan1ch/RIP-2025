package handler

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteIndexFromQuery godoc
// @Summary Удалить индекс из запроса
// @Description Удаляет связь индекса и запроса
// @Tags indexes-query
// @Produce json
// @Param index_id path int true "ID индекса"
// @Param query_id path int true "ID запроса"
// @Success 200 {object} apitypes.QueryJSON "Обновленный запрос"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes_query/{index_id}/{query_id} [delete]
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

// ChangeIndexQuery godoc
// @Summary Изменить данные индекса в запросе
// @Description Обновляет параметры индекса в конкретном запросе
// @Tags indexes-query
// @Accept json
// @Produce json
// @Param index_id path int true "ID индекса"
// @Param query_id path int true "ID запроса"
// @Param data body apitypes.IndexesQueryJSON true "Новые данные"
// @Success 200 {object} apitypes.IndexesQueryJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes_query/{index_id}/{query_id} [put]
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