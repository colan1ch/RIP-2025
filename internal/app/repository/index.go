package repository

import (
	// "database/sql"
	// "errors"
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"LAB1/internal/app/minioClient"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w:  индекс  с id %d", ErrNotFound, id)
		}
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

func (r *Repository) CreateIndex(indexJSON apitypes.IndexJSON) (ds.Index, error) {
	fmt.Print("Creating index with data:", indexJSON)
	index := apitypes.IndexFromJSON(indexJSON)
	// if index.StarRadius <= 0 {
	// 	return ds.Index{}, errors.New("неправильный радиус звезды")
	// }
	// if index.Mass <= 0 {
	// 	return ds.Index{}, errors.New("нерпавильная масса")
	// }
	err := r.db.Create(&index).First(&index).Error
	fmt.Println("Created index:", index)
	if err != nil {
		return ds.Index{}, err
	}
	return index, nil
}

func (r *Repository) ChangeIndex(id int, indexJSON apitypes.IndexJSON) (ds.Index, error) {
	index := ds.Index{}
	if id < 0 {
		return ds.Index{}, errors.New("id индекса должен быть >= 0")
	}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&index).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Index{}, fmt.Errorf("%w: индекс с id %d", ErrNotFound, id)
		}
		return ds.Index{}, err
	}
	err = r.db.Model(&index).Updates(apitypes.IndexFromJSON(indexJSON)).Error
	if err != nil {
		return ds.Index{}, err
	}
	return index, nil
}

func (r *Repository) DeleteIndex(id int) error {
	index := ds.Index{}
	if id < 0 {
		return errors.New("id должно быть >= 0")
	}

	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&index).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: индекс с id %d", ErrNotFound, id)
		}
		return err
	}
	if index.Image != "" {
		err = minio.DeleteObject(context.Background(), r.mc, minio.GetImgBucket(), index.Image)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Index{}).Where("id = ?", id).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddIndexToQuery(queryId int, indexId int) error {
	var index ds.Index
	if err := r.db.First(&index, indexId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: index with id %d", ErrNotFound, indexId)
		}
		return err
	}

	var query ds.Query
	if err := r.db.First(&query, queryId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: запрос с id %d", ErrNotFound, queryId)
		}
		return err
	}
	
	indexesQuery := ds.IndexesQuery{}
	result := r.db.Where("index_id = ? and query_id = ?", indexId, queryId).Find(&indexesQuery)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return fmt.Errorf("%w: индекс %d уже в запросе %d", ErrAlreadyExists, indexId, queryId)
	}
	return r.db.Create(&ds.IndexesQuery{
		IndexID:    uint(indexId),
		QueryID: uint(queryId),
	}).Error
}

func (r *Repository) GetModeratorAndCreatorLogin(query ds.Query) (string, string, error) {
	var creator ds.User
	var moderator ds.User

	err := r.db.Where("id = ?", query.CreatorID).First(&creator).Error
	if err != nil {
		return "", "", err
	}
	// fmt.Println(creator.Login)
	var moderatorLogin string
	if query.ModeratorID.Valid {
		err = r.db.Where("id = ?", query.ModeratorID).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) UploadImage(ctx *gin.Context, indexId int, file *multipart.FileHeader) ( ds.Index, error) {
	index_, err := r.GetIndex(indexId)
	if err != nil {
		return ds.Index{}, err
	}
	
	fileName, err := minio.UploadImage(ctx, r.mc, minio.GetImgBucket(), file, *index_)
	if err != nil {
		return ds.Index{},err
	}

	index, err := r.GetIndex(indexId)
	if err != nil {
		return ds.Index{}, err
	}
	index.Image = fileName
	err = r.db.Save(&index).Error
	if err != nil {
		return ds.Index{}, err
	}
	return *index, nil
}