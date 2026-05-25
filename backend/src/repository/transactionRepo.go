package repository

import (
	"log"
	"server/src/models"
	"server/src/types"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

func NewTransactionRepo(db *gorm.DB, rdb *redis.Client) *TransactionRepo {
	return &TransactionRepo{
		Db:  db,
		Rdb: rdb,
	}
}

func (r *TransactionRepo) GetStudentById(studentId string) (types.Student, error) {
	if len(studentId) != 24 {
		return types.Student{}, &types.RepoError{
			Message: "Invalid student ID",
			Type:    1,
		}
	}

	var student types.Student

	result := r.Db.Model(&models.Student{}).Select("id", "name", "is_blocked").First(&student, "id = ?", studentId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return types.Student{}, &types.RepoError{
				Message: "Client with " + studentId + " not found",
				Type:    1,
			}
		} else if result.Error != nil {
			log.Printf("%+v\n", result.Error)
			return types.Student{}, &types.RepoError{
				Message: "Database error!",
				Type:    0,
			}
		}
	}

	return student, nil
}

func removeDuplicateIds(ids []string) []string { // Ganti int dengan tipe data ID Anda (misal string/uint)
	keys := make(map[string]bool)
	var list []string
	for _, entry := range ids {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (r *TransactionRepo) GetItemsByIds(itemsId []string) ([]types.Item, error) {
	if len(itemsId) == 0 {
		return nil, &types.RepoError{
			Message: "Array length must not zero!",
			Type:    2,
		}
	}

	uniqueItemsId := removeDuplicateIds(itemsId)

	var item []types.Item
	err := r.Db.Model(&models.PurchaseItem{}).Select("id", "item_name", "price").Where("id IN ?", uniqueItemsId).Find(&item).Error
	if err != nil {
		return nil, &types.RepoError{
			Message: "Database error!",
			Type:    0,
		}
	}

	// Periksa apakah jumlah data yang ditemukan sama dengan jumlah ID yang dicari
	if len(item) != len(itemsId) {
		return nil, &types.RepoError{
			Message: "There's some id that not found!",
			Type:    1,
		}
	}

	return item, nil
}

func (r *TransactionRepo) CreateTransaction(studentId string, clientId string, total float64) (models.Transaction, error) {
	var newTransaction models.Transaction

	if len(studentId) != 24 {
		return models.Transaction{}, &types.RepoError{
			Message: "studentId must be a valid 24 char!",
			Type:    2,
		}
	}

	parsedClientId, err := uuid.Parse(clientId)
	if err != nil {
		return models.Transaction{}, &types.RepoError{
			Message: "Failed to parse uuid!",
			Type:    2,
		}
	}

	// isi data
	newTransaction.ID = uuid.New()
	newTransaction.StudentID = studentId
	newTransaction.ClientID = parsedClientId
	newTransaction.TotalAmount = total
	newTransaction.TransactionAt = time.Now()

	result := r.Db.Create(&newTransaction)
	if result.Error != nil {
		return models.Transaction{}, &types.RepoError{
			Message: "Database error!",
			Type:    0,
		}
	}

	return newTransaction, nil
}
