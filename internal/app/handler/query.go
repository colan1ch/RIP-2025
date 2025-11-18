package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetQueries godoc
// @Summary Получить список запросов
// @Description Возвращает запросы с возможностью фильтрации по датам и статусу
// @Tags queries
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус запроса"
// @Success 200 {array} apitypes.QueryJSON "Список запросов"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries [get]
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
	fmt.Println(queries)
	queries = h.filterQueriesByAuth(queries, ctx)




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

// GetQueryCart godoc
// @Summary Получить корзину запроса
// @Description Возвращает информацию о текущем черновике запроса пользователя
// @Tags queries
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины запроса"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/query-cart [get]
func (h *Handler) GetQueryCart(ctx *gin.Context){
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"id":          -1,
			"indexes_count": 0,
		})
		return
	}
	indexesCount := h.Repository.GetQueryCount(userID)

	if indexesCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			// "status":          "no_draft",
			"id": -1,
			"indexes_count": indexesCount,
		})
		return
	}

	query, err := h.Repository.CheckCurrentQueryDraft(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				// "status":          "no_draft",
				"id": -1,
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

// GetQuery godoc
// @Summary Получить запрос по ID
// @Description Возвращает полную информацию о запросе включая индексы
// @Tags queries
// @Produce json
// @Param id path int true "ID запроса"
// @Success 200 {object} map[string]interface{} "Данные запроса с индексами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Запрос не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/{id} [get]
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

	all := make([]interface{}, 0, len(resp)+len(resp2))
	// all = append(all, apitypes.QueryToJSON(query, creatorLogin, moderatorLogin))
	for _, p := range resp {
		all = append(all, p)
	}
	for _, pb := range resp2 {
		all = append(all, pb)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"query": apitypes.QueryToJSON(query, creatorLogin, moderatorLogin),
		"indexes": all,
		// "indexes":   resp,
		// "indexesQuery": resp2,
	})
}

// FormQuery godoc
// @Summary Сформировать запрос
// @Description Переводит запрос в статус "formed"
// @Tags queries
// @Produce json
// @Param id path int true "ID запроса"
// @Success 200 {object} apitypes.QueryJSON "Сформированный запрос"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Запрос не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/{id}/form [put]
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

// ChangeQuery godoc
// @Summary Изменить запрос
// @Description Обновляет данные запроса
// @Tags queries
// @Accept json
// @Produce json
// @Param id path int true "ID запроса"
// @Param query body apitypes.QueryJSON true "Новые данные запроса"
// @Success 200 {object} apitypes.QueryJSON "Обновленный запрос"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Запрос не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/{id}/change-query [put]
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

// DeleteQuery godoc
// @Summary Удалить запрос
// @Description Выполняет логическое удаление запроса
// @Tags queries
// @Produce json
// @Param id path int true "ID запроса"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Запрос не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/{id}/delete-query [delete]
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

// ModerateQuery godoc
// @Summary Модерировать запрос
// @Description Изменяет статус запроса (только для модераторов)
// @Tags queries
// @Accept json
// @Produce json
// @Param id path int true "ID запроса"
// @Param status body apitypes.StatusJSON true "Новый статус"
// @Success 200 {object} apitypes.QueryJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Запрос не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /queries/{id}/finish [put]
func (h *Handler) ModerateQuery(ctx *gin.Context) {
    userID, err := getUserID(ctx)
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

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

    user, err := h.Repository.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            h.errorHandler(ctx, http.StatusNotFound, err)
        } else {
            h.errorHandler(ctx, http.StatusInternalServerError, err)
        }
        return
    }
    
    if !user.IsModerator {
        h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
        return
    }

    query, err := h.Repository.ModerateQuery(id, statusJSON.Status, userID)
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

func (h *Handler) filterQueriesByAuth(queries []ds.Query, ctx *gin.Context) []ds.Query {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.Query{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return []ds.Query{}
	}
	if err != nil {
		return []ds.Query{}
	}

	if user.IsModerator {
		return queries
	}

	var userQueries []ds.Query
    for _, query := range queries {
        fmt.Println(query.ID)
        if query.CreatorID == userID {
            userQueries = append(userQueries, query)
        }
    }
    
    return userQueries

}

func (h *Handler) hasAccessToQuery(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}