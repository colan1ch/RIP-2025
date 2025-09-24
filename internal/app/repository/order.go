package repository

import (
	"LAB1/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var errNoDraft = errors.New("no draft for this user")


func (r *Repository) GetIndexesOrder(id int) ([]ds.IndexInfo, ds.Order, error) {

	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	var order ds.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return []ds.IndexInfo{}, ds.Order{}, err
	} else if creatorID != int(order.CreatorID) {
		return []ds.IndexInfo{}, ds.Order{}, errors.New("you are not allowed")
	} else if order.Status == "deleted" {
		return []ds.IndexInfo{}, ds.Order{}, errors.New("you can`t watch deleted calculations")
	}

	var indexes []ds.Index
	var indexesOrderes []ds.IndexesOrder
	sub := r.db.Table("indexes_orders").Where("order_id = ?", order.ID).Find(&indexesOrderes)
	err = r.db.Where("id IN (?)", sub.Select("index_id")).Find(&indexes).Error
	if err != nil {
		return []ds.IndexInfo{}, ds.Order{}, err
	}


	var indexesResult []ds.IndexInfo
	for _, index := range indexes {
		for _, indexesOrder := range indexesOrderes {
			if index.ID == int(indexesOrder.IndexID) {
				indexesResult = append(indexesResult, ds.IndexInfo{
					ID:                 index.ID,
					Name:              	index.Name,
					Image:            	index.Image,

					Cardinality:         index.Cardinality,
					
					RowsCount: 		indexesOrder.RowsCount,
					RecievedRows:       indexesOrder.RecievedRows,
				})
				break
			}
		}
	}

	return indexesResult, order, nil
}


func (r *Repository) CheckCurrentOrderDraft(creatorID int) (ds.Order, error) {
	var order ds.Order

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&order)
	if res.Error != nil {
		return ds.Order{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Order{}, errNoDraft
	}
	return order, nil
}


func (r *Repository) GetOrderDraft(creatorID int) (ds.Order, error) {
	order, err := r.CheckCurrentOrderDraft(creatorID)
	if err == errNoDraft {
		order = ds.Order{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&order)
		if result.Error != nil {
			return ds.Order{}, result.Error
		}
		return order, nil
	} else if err != nil {
		return ds.Order{}, err
	}
	return order, nil
}

func (r *Repository) GetOrderCount() int64 {
	var orderID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Order{}).Where("creator_id = ? AND status = ?", creatorID, "draft").Select("id").First(&orderID).Error
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.IndexesOrder{}).Where("order_id = ?", orderID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_indexes:", err)
	}

	return count
}


func (r *Repository) DeleteCalculation(orderId int) error{
	return r.db.Exec("UPDATE orders SET status = 'deleted' WHERE id = ?", orderId).Error
}