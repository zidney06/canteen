package models

import "github.com/google/uuid"

type PurchaseItem struct {
	BaseModel
	ItemName string    `json:"itemName" binding:"required" gorm:"type:varchar(100);not null"`
	Price    float64   `json:"price" binding:"required" gorm:"type:numeric;not null"`
	ClientId uuid.UUID `json:"clientId" gorm:"type:uuid;not null"`
}
