package repository

import (
	// "database/sql"
	// "errors"
	"fmt"

	"LAB1/internal/app/ds"

)

func (r *Repository) GetIndexes() ([]ds.Index, error) {
	var indexes []ds.Index
	err := r.db.Order("id").Where("is_delete = false").Find(&indexes).Error
	if err != nil {
		return nil, err
	}
	if len(indexes) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return indexes, nil
}

func (r *Repository) GetIndex(id int) (*ds.Index, error) {
	index := ds.Index{}
	err := r.db.Order("id").Where("id = ? and is_delete = ?", id, false).First(&index).Error
	if err != nil {
		return &ds.Index{}, err
	}
	return &index, nil
}

func (r *Repository) GetIndexesByName(name string) ([]ds.Index, error) {
	var indexes []ds.Index
	err := r.db.Order("id").Where("name ILIKE ? and is_delete = ?", "%"+name+"%", false).Find(&indexes).Error
	if err != nil {
		return nil, err
	}
	return indexes, nil
}


func (r *Repository) AddIndexToQuery(queryId int, indexId int) error {
	var index ds.Index
	if err := r.db.First(&index, indexId).Error; err != nil {
		return err
	}

	var query ds.Query
	if err := r.db.First(&query, queryId).Error; err != nil {
		return err
	}
	indexesQuery := ds.IndexesQuery{}
	result := r.db.Where("index_id = ? and query_id = ?", indexId, queryId).Find(&indexesQuery)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.IndexesQuery{
		IndexID:    uint(indexId),
		QueryID: uint(queryId),
	}).Error
}