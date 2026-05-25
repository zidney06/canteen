package repository

import (
	"errors"
	"log"
	"server/src/models"
	"server/src/types"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

func NewProductRepo(db *gorm.DB, rdb *redis.Client) *ProductRepo {
	return &ProductRepo{
		Db:  db,
		Rdb: rdb,
	}
}

func (r *ProductRepo) GetProductList(clientId string) ([]types.Product, error) {
	var products []types.Product

	result := r.Db.Model(&models.PurchaseItem{}).Select("item_name", "id", "price").Where("client_id = ?", clientId).Find(&products)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, &types.RepoError{
			Message: "Database error!",
			Type:    0,
		}
	}

	if len(products) == 0 {
		return nil, &types.RepoError{
			Message: "Product with " + clientId + " not found!",
			Type:    1,
		}
	}

	return products, nil
}

func (r *ProductRepo) AddProducts(items []types.Product, clientId string) (int, error) {
	if len(items) == 0 {
		return 0, &types.ServiceError{
			Message:    "Items length must not zero!",
			HttpStatus: 400,
		}
	}

	var products []models.PurchaseItem
	parsedUuid, err := uuid.Parse(clientId)
	if err != nil {
		log.Println(err.Error())
		return 0, &types.RepoError{
			Message: "Internal server error",
		}
	}

	// lakukan remapping pada products
	for _, item := range items {
		p := models.PurchaseItem{
			ItemName: item.ItemName,
			Price:    item.Price,
			ClientId: parsedUuid,
		}

		p.ID = uuid.New()

		products = append(products, p)
	}

	result := r.Db.Create(&products)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, &types.RepoError{
			Message: "Failed to save to database",
		}
	}

	return int(result.RowsAffected), nil
}

func (r *ProductRepo) UpdateProduct(productId string, itemName string, price float64) (types.Product, error) {
	var updatedProduct types.Product

	result := r.Db.Model(&models.PurchaseItem{}).Clauses(clause.Returning{}).Where("id = ?", productId).Updates(map[string]any{"item_name": itemName, "price": price}).Scan(&updatedProduct)
	if result.Error != nil {
		log.Println(result.Error)
		return types.Product{}, &types.RepoError{
			Message: "Internal server error!",
			Type:    0,
		}
	} else if result.RowsAffected == 0 {
		return types.Product{}, &types.RepoError{
			Message: "Product with " + productId + " not found!",
			Type:    1,
		}
	}

	return updatedProduct, nil
}

func (r *ProductRepo) DeleteProduct(productIds []string) error {
	var arrayString []string
	var repoErr *types.RepoError

	if len(productIds) == 0 {
		return &types.RepoError{
			Message: "array of productIds must not zero!",
			Type:    2,
		}
	}

	for _, id := range productIds {
		if len(strings.TrimSpace(id)) > 0 {
			arrayString = append(arrayString, strings.TrimSpace(id))
		}
	}

	if len(arrayString) != len(productIds) {
		return &types.RepoError{
			Message: `Array element must not empty`,
			Type:    2,
		}
	}

	// transaction. jika ada id yang tidak ada di db batalkan proses
	err := r.Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&models.PurchaseItem{}, arrayString)
		if result.Error != nil {
			log.Println("Failed:", result.Error)
			return &types.RepoError{
				Message: "Internal server error",
				Type:    0,
			}
		}
		if int(result.RowsAffected) < len(productIds) {
			return &types.RepoError{
				Message: "There is id that not exits in the database!",
				Type:    2,
			}
		}

		return nil
	})

	if err != nil && errors.As(err, &repoErr) {
		return repoErr
	}

	return nil
}
