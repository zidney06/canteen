package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	// tidak menggunakan gorm.Model karena disini butuhnya adalah uuid sedangkan gorm.Model menggunakan uint
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
