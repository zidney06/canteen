package models

import (
	"time"

	"github.com/google/uuid"
)

// ini buat dikirimkan ke client saat transakasi berhasil agar digunakan sebagai struk
type Transaction struct {
	BaseModel
	StudentID     string    `json:"studentId" gorm:"type:varchar(24);index;not null"`
	Student       Student   `json:"student" gorm:"foreignKey:StudentID"`
	ClientID      uuid.UUID `json:"clientId" gorm:"type:uuid;index;not null"`
	Client        Client    `json:"client" gorm:"foreignKey:ClientID"`
	TotalAmount   float64   `json:"totalAmount" gorm:"type:decimal(10,2);not null"`
	TransactionAt time.Time `json:"transactionAt" gorm:"index"`
}
