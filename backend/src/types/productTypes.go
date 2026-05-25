package types

import "github.com/google/uuid"

type Product struct {
	Id       uuid.UUID `json:"id"`
	ItemName string    `json:"itemName"`
	Price    float64   `json:"price"`
}

type AddProductsBody struct {
	Products []Product
}

type UpdateProductBody struct {
	ItemName string  `json:"itemName" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}

type DeleteProductsBody struct {
	ProductsIds []string `json:"productIds" binding:"required"`
}
