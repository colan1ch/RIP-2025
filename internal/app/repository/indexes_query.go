package repository

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (r *Repository) DeleteIndexFromQuery(queryId int, indexId int) (ds.Query, error) {
	// userId := r.userId
    // if userId == 0 {
    //     return ds.Query{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Query{}, err
	// }
    
	var query ds.Query
	err := r.db.Where("id = ?", queryId).First(&query).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Query{}, fmt.Errorf("%w: запрос с id %d", ErrNotFound, queryId)
		}
		return ds.Query{}, err
	}
    
	// if query.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Query{}, fmt.Errorf("%w: Вы не создатель этого исследования", ErrNotAllowed)
	// }
    
	err = r.db.Where("index_id = ? and Query_id = ?", indexId, queryId).Delete(&ds.IndexesQuery{}).Error
	if err != nil {
		return ds.Query{}, err
	}
	return query, nil
}

func (r *Repository) ChangeIndexQuery(queryId int, indexId int, indexesQueryJSON apitypes.IndexesQueryJSON) (ds.IndexesQuery, error) {
	var indexesQuery ds.IndexesQuery
	err := r.db.Model(&indexesQuery).Where("index_id = ? and query_id = ?", indexId, queryId).Updates(apitypes.IndexesQueryFromJSON(indexesQueryJSON)).First(&indexesQuery).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.IndexesQuery{}, fmt.Errorf("%w: индексы в запросе", ErrNotFound)
		}
		return ds.IndexesQuery{}, err
	}
	return indexesQuery, nil
}