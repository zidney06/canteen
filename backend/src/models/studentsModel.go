package models

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID        string         `json:"id" gorm:"primaryKey; size:24"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	IsBlocked bool           `json:"isBlocked" gorm:"default:false"`
}
