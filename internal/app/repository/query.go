package repository

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetQueries(from, to time.Time, status string) ([]ds.Query, error) {
	var queries []ds.Query
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("date_create > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_create < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("id").Find(&queries).Error
	if err != nil {
		return nil, err
	}
	return queries, nil
}

func (r *Repository) GetIndexesQueries(queryId int) ([]ds.IndexesQuery, error) {
	var indexesQueries []ds.IndexesQuery
	err := r.db.Where("query_id = ?", queryId).Find(&indexesQueries).Error
	if err != nil {
		return nil, err
	}
	return indexesQueries, nil
}

func (r *Repository) GetIndexesQuery(indexId int, queryId int) (ds.IndexesQuery, error) {
	var indexesQuery ds.IndexesQuery
	err := r.db.Where("index_id = ? and query_id = ?", indexId, queryId).First(&indexesQuery).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.IndexesQuery{}, fmt.Errorf("%w: indexes query not found", ErrNotFound)
		}
		return ds.IndexesQuery{}, err
	}
	return indexesQuery, nil
}

func (r *Repository) GetQueryIndexes(id int) ([]ds.Index, ds.Query, error) {
	query, err := r.GetSingleQuery(id)
	if err != nil {
		return []ds.Index{}, ds.Query{}, err
	}

	var indexes []ds.Index
	sub := r.db.Table("indexes_queries").Where("Query_id = ?", query.ID)
	err = r.db.Order("id DESC").Where("id IN (?)", sub.Select("index_id")).Find(&indexes).Error

	if err != nil {
		return []ds.Index{}, ds.Query{}, err
	}

	return indexes, query, nil
}

func (r *Repository) CheckCurrentQueryDraft(creatorID uuid.UUID) (ds.Query, error) {
	// if creatorID == 0 {
	//     return ds.Query{}, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
	// }

	var query ds.Query
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&query)
	if res.Error != nil {
		return ds.Query{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Query{}, ErrNoDraft
	}
	return query, nil
}

func (r *Repository) GetQueryDraft(creatorID uuid.UUID) (ds.Query, bool, error) {
	// if creatorID == 0 {
	//     return ds.Query{}, false, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
	// }

	query, err := r.CheckCurrentQueryDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
		query = ds.Query{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&query)
		if result.Error != nil {
			return ds.Query{}, false, result.Error
		}
		return query, true, nil
	} else if err != nil {
		return ds.Query{}, false, err
	}
	return query, true, nil
}

func (r *Repository) GetQueryCount(creatorID uuid.UUID) int64 {

	var count int64
	query, err := r.CheckCurrentQueryDraft(creatorID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.IndexesQuery{}).Where("query_id = ?", query.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_indexes:", err)
	}

	return count
}

func (r *Repository) DeleteCalculation(queryId int) error {
	return r.db.Exec("UPDATE queries SET status = 'deleted' WHERE id = ?", queryId).Error
}

func (r *Repository) GetSingleQuery(id int) (ds.Query, error) {
	if id < 0 {
		return ds.Query{}, errors.New("неверное id, должно быть >= 0")
	}

	// userId := r.GetUserID()
	// if userId == 0 {
	//     return ds.Query{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
	// }

	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Query{}, err
	// }

	var query ds.Query
	err := r.db.Where("id = ?", id).First(&query).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Query{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Query{}, err
	} else if query.Status == "deleted" {
		return ds.Query{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return query, nil
}

func (r *Repository) FormQuery(queryId int, status string) (ds.Query, error) {
	query, err := r.GetSingleQuery(queryId)
	if err != nil {
		return ds.Query{}, err
	}

	if query.Status != "draft" {
		return ds.Query{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	if status != "deleted" {
		if query.DateQuery == "" {
			return ds.Query{}, errors.New("вы не написали дату запроса")
		}
		indexesQuery, _ := r.GetIndexesQueries(query.ID)
		for _, indexQuery := range indexesQuery {
			if indexQuery.RowsCount <= 0 {
				return ds.Query{}, errors.New("вы не написали количество строк в таблице")
			}
			if indexQuery.Cardinality <= 0 {
				return ds.Query{}, errors.New("вы не написали кардинальность индекса")
			}
		}
	}

	err = r.db.Model(&query).Updates(ds.Query{
		Status: status,
		DateForm: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Query{}, err
	}

	return query, nil
}

func (r *Repository) ChangeQuery(id int, queryJSON apitypes.QueryJSON) (ds.Query, error) {
	query := ds.Query{}
	if id < 0 {
		return ds.Query{}, errors.New("неправильное id, должно быть >= 0")
	}
	if queryJSON.DateQuery == "" {
		return ds.Query{}, errors.New("неправильная дата запроса")
	}
	err := r.db.Where("id = ? and status != 'deleted'", id).First(&query).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Query{}, fmt.Errorf("%w: запрос с id %d", ErrNotFound, id)
		}
		return ds.Query{}, err
	}
	err = r.db.Model(&query).Updates(apitypes.QueryFromJSON(queryJSON)).Error
	if err != nil {
		return ds.Query{}, err
	}
	return query, nil
}

func CalculateRecievedRows(cardinality int, dateQuery string, rows_count int) (int, error) {
	if dateQuery == "" {
		return 0, errors.New("неправильная дата запроса")
	}
	// if indexShine < 0 || indexShine > 7 {
	// 	return 0, errors.New("неправильный блеск")
	// }
	return rows_count - cardinality, nil
}

func (r *Repository) ModerateQuery(id int, status string, currUserId uuid.UUID) (ds.Query, error) {
	if status != "completed" && status != "rejected" {
		return ds.Query{}, errors.New("неверный статус")
	}

	query, err := r.GetSingleQuery(id)
	if err != nil {
		return ds.Query{}, err
	} else if query.Status != "formed" {
		return ds.Query{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&query).Updates(ds.Query{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: uuid.NullUUID{
			UUID:  currUserId,
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Query{}, err
	}

	// if status == "completed" {
		// indexesQuery, err := r.GetIndexesQueries(query.ID)
		// if err != nil {
		// 	return ds.Query{}, err
		// }
		// executionTime будет рассчитан Django async сервисом и отправлен обратно
		// Здесь мы просто рассчитываем recieved_rows для каждого индекса
		// for _, indexQuery := range indexesQuery {
		// 	recievedRows, err := CalculateRecievedRows(indexQuery.Cardinality, query.DateQuery, indexQuery.RowsCount)
		// 	if err != nil {
		// 		return ds.Query{}, err
		// 	}
		// 	err = r.db.Model(&indexQuery).Updates(ds.IndexesQuery{
		// 		RecievedRows: int(recievedRows),
		// 	}).Error
		// 	if err != nil {
		// 		return ds.Query{}, err
		// 	}
		// }
	// }
	return query, nil
}

func CalculateExecutionTime(indexesQuery []ds.IndexesQuery) (float64, error) {
	sum := 0.0
	pr := 1.0
	for _, indexQuery := range indexesQuery {
		sum += math.Log2(float64(indexQuery.Cardinality))
		pr *= float64(indexQuery.RowsCount) / float64(indexQuery.Cardinality)
	}
	return float64(sum + pr), nil
}

func (r *Repository) UpdateQueryStatus(id int, status string) (ds.Query, error) {
	var query ds.Query
	if err := r.db.Model(&query).Where("id = ?", id).Update("status", status).Error; err != nil {
		return ds.Query{}, err
	}
	return query, nil
}

// GetQueryWithIndexes получает Query с preload связанных IndexesQueries
func (r *Repository) GetQueryWithIndexes(queryID int) (ds.Query, error) {
	var query ds.Query
	if err := r.db.Preload("IndexesQueries").First(&query, queryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Query{}, fmt.Errorf("%w: query not found", ErrNotFound)
		}
		return ds.Query{}, err
	}
	return query, nil
}

// UpdateQueryExecutionTime обновляет время выполнения запроса
func (r *Repository) UpdateQueryExecutionTime(queryID int, executionTime int) error {
	if err := r.db.Model(&ds.Query{}).Where("id = ?", queryID).Update("execution_time", executionTime).Error; err != nil {
		return err
	}
	return nil
}

// UpdateIndexQueryReceivedRows обновляет количество полученных строк для IndexesQuery
func (r *Repository) UpdateIndexQueryReceivedRows(indexID int, receivedRows int) error {
	if err := r.db.Model(&ds.IndexesQuery{}).Where("id = ?", indexID).Update("recieved_rows", receivedRows).Error; err != nil {
		return err
	}
	return nil
}

// GetCompletedIndexesCount возвращает количество обработанных индексов
func (r *Repository) GetCompletedIndexesCount(queryID int) (int, error) {
	var count int64
	if err := r.db.Model(&ds.IndexesQuery{}).
		Where("query_id = ? AND recieved_rows > 0", queryID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
