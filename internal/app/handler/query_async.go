package handler

import (
	"bytes"
	"encoding/json"
	"io"

	// "log"
	"net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const ASYNC_SERVICE_URL = "http://localhost:8000/api/calculate/"
const SECRET_KEY = "SuperSecretKey"

type AsyncIndexPayload struct {
	IndexID        int    `json:"index_id"`
	Cardinality    int    `json:"cardinality"`
	RowsCount      int    `json:"rows_count"`
	DateQuery      string `json:"date_query"`
	ReceivedRowsID int    `json:"received_rows_id"`
}

type AsyncQueryPayload struct {
	QueryID     int                 `json:"query_id"`
	IndexesData []AsyncIndexPayload `json:"indexes_data"`
}

type AsyncResultIndex struct {
	IndexID      int    `json:"index_id"`
	ReceivedRows int    `json:"received_rows"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

type AsyncResultPayload struct {
	QueryID       int                `json:"query_id"`
	ExecutionTime float64            `json:"execution_time"`
	IndexResults  []AsyncResultIndex `json:"index_results"`
	Timestamp     float64            `json:"timestamp"`
}

// TriggerAsyncCalculation отправляет запрос на async сервис
func (h *Handler) TriggerAsyncCalculation(queryID int) {
	go func() {
		// Получаем заявку с индексами
		query, err := h.Repository.GetQueryWithIndexes(queryID)
		if err != nil {
			logrus.Errorf("Error fetching query %d: %v", queryID, err)
			return
		}

		// Подготавливаем данные индексов
		indexesData := make([]AsyncIndexPayload, 0, len(query.IndexesQueries))
		for _, idx := range query.IndexesQueries {
			indexesData = append(indexesData, AsyncIndexPayload{
				IndexID:        int(idx.IndexID),
				Cardinality:    idx.Cardinality,
				RowsCount:      idx.RowsCount,
				DateQuery:      query.DateQuery,
				ReceivedRowsID: int(idx.ID),
			})
		}

		payload := AsyncQueryPayload{
			QueryID:     queryID,
			IndexesData: indexesData,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			logrus.Errorf("Error marshaling payload for query ID %d: %v", queryID, err)
			return
		}

		resp, err := http.Post(ASYNC_SERVICE_URL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			logrus.Errorf("Error calling async service for query ID %d: %v", queryID, err)
			return
		}
		defer resp.Body.Close()

		// Полностью читаем ответ перед закрытием соединения
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Error reading response body for query ID %d: %v", queryID, err)
			return
		}

		if resp.StatusCode != http.StatusAccepted {
			logrus.Errorf("Async service returned non-202 status for query ID %d: %s. Response: %s", queryID, resp.Status, string(respBody))
			return
		}

		logrus.Infof("Successfully sent task to async service for query ID %d. Response: %s", queryID, string(respBody))
	}()
}

// UpdateQueryResult получает результаты от async сервиса
func (h *Handler) UpdateQueryResult(c *gin.Context) {
	var payload AsyncResultPayload
	if err := c.BindJSON(&payload); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	// Проверка авторизации через Bearer токен
	authHeader := c.GetHeader("Authorization")
	expectedAuth := "Bearer " + SECRET_KEY
	if authHeader != expectedAuth {
		h.errorHandler(c, http.StatusForbidden, nil)
		return
	}

	logrus.Infof("Received calculation results for query ID %d with execution time %.2f seconds", payload.QueryID, payload.ExecutionTime)

	// Обновляем execution_time в Query
	if err := h.Repository.UpdateQueryExecutionTime(payload.QueryID, int(payload.ExecutionTime)); err != nil {
		logrus.Errorf("Failed to update execution_time for query %d: %v", payload.QueryID, err)
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}
	logrus.Infof("Updated execution_time for query %d to %d", payload.QueryID, int(payload.ExecutionTime))

	// Обновляем received_rows для каждого индекса
	successCount := 0
	for _, result := range payload.IndexResults {
		if result.Success {
			if err := h.Repository.UpdateIndexQueryReceivedRows(result.IndexID, result.ReceivedRows); err != nil {
				logrus.Errorf("Failed to update recieved_rows for index %d: %v", result.IndexID, err)
			} else {
				logrus.Infof("Updated index %d received_rows to %d", result.IndexID, result.ReceivedRows)
				successCount++
			}
		}
	}

	logrus.Infof("Successfully updated results for query %d: %d indexes updated", payload.QueryID, successCount)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "updated_indexes": successCount})
}

// GetCompletedIndexesCount возвращает количество обработанных индексов
func (h *Handler) GetCompletedIndexesCount(queryID int) (int, error) {
	return h.Repository.GetCompletedIndexesCount(queryID)
}
