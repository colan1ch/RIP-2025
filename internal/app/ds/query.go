package ds

import (
	"database/sql"
	"time"
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
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	DateQuery time.Time
	Status       string `gorm:"type:varchar(15);not null;check:status IN ('draft','deleted','formed','completed','rejected')"`
	DateCreate   time.Time `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID    int       `gorm:"not null"`
	ModeratorID sql.NullInt64 `gorm:"default:null"`

	Creator User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
	ExecutionTime int
}
