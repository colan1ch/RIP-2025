package ds

import (
	"database/sql"
	"time"
)

type OrderStatus string

const (
  OrderStatusDraft     OrderStatus = "draft"
  OrderStatusDeleted   OrderStatus = "deleted"
  OrderStatusFormed    OrderStatus = "formed"
  OrderStatusCompleted OrderStatus = "completed"
  OrderStatusRejected  OrderStatus = "rejected"
)

type Order struct {
	ID           uint   `gorm:"primaryKey"`
	DateOrder time.Time
	Status       string `gorm:"type:varchar(15);not null;check:status IN ('draft','deleted','formed','completed','rejected')"`
	DateCreate   time.Time `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID    int       `gorm:"not null"`
	ModeratorID sql.NullInt64 `gorm:"default:null"`

	Creator Users `gorm:"foreignKey:CreatorID"`
	Moderator Users `gorm:"foreignKey:ModeratorID"`
	ExecutionTime int
}
