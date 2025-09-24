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
// 	query := "SELECT id, image, name, description, distance, mass, discovery, star_radius FROM indexes WHERE id = $1"
// 	row := r.db.Raw(query, id).Row()
// 	index := &ds.Index{}

//    err := row.Scan(
// 		&index.ID,
//       &index.Image,
//       &index.Name,
//       &index.Description,
//       &index.Distance,
//       &index.Mass,
// 	  &index.Discovery,
//       &index.StarRadius,
//    )
//    if err != nil {
//       if errors.Is(err, sql.ErrNoRows) {
//          return nil, nil // Возвращаем nil, если записи нет
//       }
//       return nil, err
//    }
// 	return index, nil

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


func (r *Repository) AddIndexToOrder(orderId int, indexId int) error {
	var index ds.Index
	if err := r.db.First(&index, indexId).Error; err != nil {
		return err
	}

	var order ds.Order
	if err := r.db.First(&order, orderId).Error; err != nil {
		return err
	}
	indexesOrder := ds.IndexesOrder{}
	result := r.db.Where("index_id = ? and order_id = ?", indexId, orderId).Find(&indexesOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.IndexesOrder{
		IndexID:    uint(indexId),
		OrderID: uint(orderId),
	}).Error
}