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

// GetIndexes godoc
// @Summary Получить список индексов
// @Description Возвращает все индексы или фильтрует по названию
// @Tags indexes
// @Produce json
// @Param index_name query string false "Название индекса для поиска"
// @Success 200 {array} apitypes.IndexJSON "Список индексов"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /indexes [get]
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

// GetIndex godoc
// @Summary Получить индекс по ID
// @Description Возвращает информацию об индексе по её идентификатору
// @Tags indexes
// @Produce json
// @Param id path int true "ID индекса"
// @Success 200 {object} apitypes.IndexJSON "Данные индекса"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Индекс не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /indexes/{id} [get]
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

// CreateIndex godoc
// @Summary Создать новый индекс
// @Description Создает новый индекс и возвращает его данные
// @Tags indexes
// @Accept json
// @Produce json
// @Param index body apitypes.IndexJSON true "Данные нового индекса"
// @Success 201 {object} apitypes.IndexJSON "Созданный индекс"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes/create-index [post]
func (h *Handler) CreateIndex(ctx *gin.Context) {
	fmt.Println("Creating index with data:")

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

// DeleteIndex godoc
// @Summary Удалить индекс
// @Description Выполняет логическое удаление индекса по ID
// @Tags indexes
// @Produce json
// @Param id path int true "ID индекса"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Индекс не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes/{id}/delete-index [delete]
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

	ctx.JSON(http.StatusOK, gin.H{
		// "status": "deleted",
	})
}

// ChangeIndex godoc
// @Summary Изменить данные индекса
// @Description Обновляет информацию об индексе по ID
// @Tags indexes
// @Accept json
// @Produce json
// @Param id path int true "ID индекса"
// @Param index body apitypes.IndexJSON true "Новые данные индекса"
// @Success 200 {object} apitypes.IndexJSON "Обновленный индекс"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Индекс не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes/{id}/change-index [put]
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

// AddIndexToQuery godoc
// @Summary Добавить индекс в запрос
// @Description Добавляет индекс в черновик запроса пользователя
// @Tags indexes
// @Produce json
// @Param id path int true "ID индекса"
// @Success 200 {object} apitypes.QueryJSON "Запрос с добавленным индексом"
// @Success 201 {object} apitypes.QueryJSON "Создан новый запрос"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Индекс не найден"
// @Failure 409 {object} map[string]string "Индекс уже в запросе"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes/{id}/add-to-query [post]
func (h *Handler) AddIndexToQuery(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	query, created, err := h.Repository.GetQueryDraft(userID)
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

// UploadImage godoc
// @Summary Загрузить изображение для индекса
// @Description Загружает изображение для индекса и возвращает обновленные данные
// @Tags indexes
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID индекса"
// @Param image formData file true "Изображение индекса"
// @Success 200 {object} map[string]interface{} "Статус загрузки и данные индекса"
// @Failure 400 {object} map[string]string "Неверный запрос или файл"
// @Failure 404 {object} map[string]string "Индекс не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /indexes/{id}/create-image [post]
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