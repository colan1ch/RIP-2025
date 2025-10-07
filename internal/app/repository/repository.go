package repository

import (
	"errors"

	"github.com/minio/minio-go/v7"
	minioClient "LAB1/internal/app/minioClient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Добавляем кастомные ошибки
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotAllowed    = errors.New("not allowed")
	ErrNoDraft       = errors.New("no draft for this user")
)

type Repository struct {
	db *gorm.DB
	mc     *minio.Client
	userId int
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	mc, err := minioClient.InitMinio()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		mc: mc,
		userId: 0,
	}, nil
}

func (r *Repository) GetUserID() (int) {
	return r.userId
}

func (r *Repository) SetUserID(id int) {
	r.userId = id
}

func (r *Repository) SignOut() {
	r.userId = 0
}