package types

import "time"

type ItemType struct {
	Id       string `json:"id" binding:"required"`
	Quantity uint8  `json:"quantity" binding:"required"`
}

type MakeTransaactionBody struct {
	EncryptedId string     `json:"encryptedId" binding:"required"`
	Items       []ItemType `json:"items" binding:"required"`
}

type Item struct {
	ID       string
	ItemName string
	Price    float64
}

type BuyedItems struct {
	Id       string
	ItemName string
	Price    float64
	Quantity uint
}

type TransactionResult struct {
	BuyerName     string
	BuyedItems    []BuyedItems
	TotalAmount   uint
	CashierName   string
	TransactionAt time.Time
}

type Student struct {
	ID        string
	Name      string
	IsBlocked bool
}
