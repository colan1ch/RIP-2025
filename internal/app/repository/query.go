package repository

import (
	"LAB1/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var errNoDraft = errors.New("no draft for this user")


func (r *Repository) GetIndexesQuery(id int) ([]ds.IndexInfo, ds.Query, error) {

	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	var query ds.Query
	err := r.db.Where("id = ?", id).First(&query).Error
	if err != nil {
		return []ds.IndexInfo{}, ds.Query{}, err
	} else if creatorID != int(query.CreatorID) {
		return []ds.IndexInfo{}, ds.Query{}, errors.New("you are not allowed")
	} else if query.Status == "deleted" {
		return []ds.IndexInfo{}, ds.Query{}, errors.New("you can`t watch deleted queries")
	}

	var indexes []ds.Index
	var indexesQueryes []ds.IndexesQuery
	sub := r.db.Table("indexes_queries").Where("query_id = ?", query.ID).Find(&indexesQueryes).Order("id")
	err = r.db.Where("id IN (?)", sub.Select("index_id")).Find(&indexes).Error
	if err != nil {
		return []ds.IndexInfo{}, ds.Query{}, err
	}


	var indexesResult []ds.IndexInfo
	var position int = 0
	for _, index := range indexes {
		for _, indexesQuery := range indexesQueryes {
			if index.ID == int(indexesQuery.IndexID) {
				position += 1
				indexesResult = append(indexesResult, ds.IndexInfo{
					ID:                 index.ID,
					Name:              	index.Name,
					Image:            	index.Image,

					Cardinality:         index.Cardinality,
					
					RowsCount: 		indexesQuery.RowsCount,
					RecievedRows:       indexesQuery.RecievedRows,
					PositionInQuery: position,
					TableField: indexesQuery.TableField,
				})
				break
			}
		}
	}

	return indexesResult, query, nil
}


func (r *Repository) CheckCurrentQueryDraft(creatorID int) (ds.Query, error) {
	var query ds.Query

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&query)
	if res.Error != nil {
		return ds.Query{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Query{}, errNoDraft
	}
	return query, nil
}


func (r *Repository) GetQueryDraft(creatorID int) (ds.Query, error) {
	query, err := r.CheckCurrentQueryDraft(creatorID)
	if err == errNoDraft {
		query = ds.Query{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&query)
		if result.Error != nil {
			return ds.Query{}, result.Error
		}
		return query, nil
	} else if err != nil {
		return ds.Query{}, err
	}
	return query, nil
}

func (r *Repository) GetQueryCount() int64 {
	var queryID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Query{}).Where("creator_id = ? AND status = ?", creatorID, "draft").Select("id").First(&queryID).Error
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.IndexesQuery{}).Where("query_id = ?", queryID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_indexes:", err)
	}

	return count
}


func (r *Repository) DeleteQuery(queryId int) error{
	return r.db.Exec("UPDATE queries SET status = 'deleted' WHERE id = ?", queryId).Error
}

func (r *Repository) UpdateRowsCount(queryId int, indexId int, rowsCount int) error{
	return r.db.Exec("UPDATE indexes_queries SET rows_count = ? WHERE query_id = ? AND index_id = ?", rowsCount, queryId, indexId).Error
}

func (r *Repository) UpdateRecievedRows(queryId int, indexId int, recievedRows int) error{
	return r.db.Exec("UPDATE indexes_queries SET recieved_rows = ? WHERE query_id = ? AND index_id = ?", recievedRows, queryId, indexId).Error
}

func (r *Repository) UpdateTableField(queryId int, indexId int, tableField string) error{
	return r.db.Exec("UPDATE indexes_queries SET table_field = ? WHERE query_id = ? AND index_id = ?", tableField, queryId, indexId).Error
}