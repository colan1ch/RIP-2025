package ds

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type QueryStatus string

const (
	QueryStatusDraft     QueryStatus = "draft"
	QueryStatusDeleted   QueryStatus = "deleted"
	QueryStatusFormed    QueryStatus = "formed"
	QueryStatusCompleted QueryStatus = "completed"
	QueryStatusRejected  QueryStatus = "rejected"
)

type Query struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	DateQuery   string
	Status      string        `gorm:"type:varchar(15);not null;check:status IN ('draft','deleted','formed','completed','rejected')"`
	DateCreate  time.Time     `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID   uuid.UUID     `gorm:"not null"`
	ModeratorID uuid.NullUUID `gorm:"default:null"`

	Creator        User           `gorm:"foreignKey:CreatorID"`
	Moderator      User           `gorm:"foreignKey:ModeratorID"`
	IndexesQueries []IndexesQuery `gorm:"foreignKey:QueryID"`
	ExecutionTime  int
}
