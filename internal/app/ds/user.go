package ds

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID   `gorm:"primaryKey;autoIncrement"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}
