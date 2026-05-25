package service

import (
	"errors"
	"math"
	"os"

	// "math"
	"server/src/repository"
	"server/src/types"
	"server/src/utils"
	// "server/src/utils"
)

type TransactionService struct {
	Repo        *repository.TransactionRepo
	StudentRepo *repository.StudentRepo
	ClientRepo  *repository.ClientRepo
}

func NewTransactionService(repo *repository.TransactionRepo, studentRepo *repository.StudentRepo, clientRepo *repository.ClientRepo) *TransactionService {
	return &TransactionService{
		Repo:        repo,
		StudentRepo: studentRepo,
		ClientRepo:  clientRepo,
	}
}

func (s *TransactionService) GetOrCreateStudent(studentId string) (types.Student, error) {
	if len(studentId) != 24 {
		return types.Student{}, &types.ServiceError{
			Message:    "Invalid student ID",
			HttpStatus: 400,
		}
	}

	student, err := s.Repo.GetStudentById(studentId)
	var repoErr *types.RepoError
	if err != nil && errors.As(err, &repoErr) {
		switch repoErr.Type {
		// jika data tidak ada di db A
		case 1:
			// buat data student baru di db A
			newStudent, createErr := s.StudentRepo.CreateNewStudent(studentId)
			if createErr != nil && errors.As(createErr, &repoErr) {
				switch repoErr.Type {
				case 1:
					return types.Student{}, &types.ServiceError{
						Message:    repoErr.Message,
						HttpStatus: 404,
					}
				case 4:
					return types.Student{}, &types.ServiceError{
						Message:    repoErr.Message,
						HttpStatus: 401,
					}
				default:
					return types.Student{}, &types.ServiceError{
						Message:    "Failed to create new student",
						HttpStatus: 500,
					}
				}
			}

			return newStudent, nil
		case 2:
			return types.Student{}, &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 400,
			}
		default:
			return types.Student{}, &types.ServiceError{
				Message:    "Internal server error!",
				HttpStatus: 500,
			}
		}
	}

	// kalau data ada di db A, kembalikan data tersebut
	return student, nil
}

func (s *TransactionService) MakeTransaction(encriptedStudentId string, items []types.ItemType, clientId string, clientName string) (types.TransactionResult, error) {
	var repoErr *types.RepoError
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	decryptedText, err := utils.Decrypt(encriptedStudentId, secretKey)
	if err != nil {
		return types.TransactionResult{}, &types.ServiceError{
			Message:    "Failed to decrypt student id",
			HttpStatus: 500,
		}
	}

	// cek ke db apk A apakah id siswa sudah ada
	student, err := s.GetOrCreateStudent(decryptedText)
	if err != nil {
		return types.TransactionResult{}, err
	}

	if student.IsBlocked {
		return types.TransactionResult{}, &types.ServiceError{
			Message:    "Student blocked!",
			HttpStatus: 400,
		}
	}

	// cek apakah id yang ada dalam items ada dalam db apk A
	var itemIds []string

	for _, i := range items {
		itemIds = append(itemIds, i.Id)
	}

	// untuk memastikan kalau id barang yang di kirim valid.
	data, err := s.Repo.GetItemsByIds(itemIds)
	if err != nil {
		return types.TransactionResult{}, utils.RepoToServiceError(err)
	}

	// hitung yang harus di bayar
	var total int64
	qtyMap := make(map[string]int)
	for _, item := range items {
		qtyMap[item.Id] = int(item.Quantity)
	}

	for _, product := range data {
		qty := qtyMap[product.ID]

		subTotal := int64(math.Round(product.Price)) * int64(qty)
		total += subTotal
	}

	mapA := make(map[string]int)
	for _, itemA := range items {
		mapA[itemA.Id] = int(itemA.Quantity)
	}

	var buyedItems []types.BuyedItems

	for _, itemB := range data {
		quantity := mapA[itemB.ID]

		merge := types.BuyedItems{
			Id:       itemB.ID,
			ItemName: itemB.ItemName,
			Price:    itemB.Price,
			Quantity: uint(quantity),
		}

		buyedItems = append(buyedItems, merge)
	}

	// simpan hasil transaksi ke db
	trs, err := s.Repo.CreateTransaction(decryptedText, clientId, float64(total))

	// ini bisa di buat menjadi function helper. jadi kita membuat sebuah function yang parameternya adalah err
	// dengan type repoError. lakukan pengecekan err != nil && errors.As(err, &repoErr) lalu lakukan switch case
	// yang mereturn bentuk errornya.
	if err != nil && errors.As(err, &repoErr) {
		switch repoErr.Type {
		case 2:
			return types.TransactionResult{}, &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 400,
			}
		default:
			return types.TransactionResult{}, &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 500,
			}
		}
	}

	// jumlah barang yang dibeli belum di kirim
	return types.TransactionResult{
		BuyerName:     student.Name,
		BuyedItems:    buyedItems,
		TotalAmount:   uint(trs.TotalAmount),
		CashierName:   clientName,
		TransactionAt: trs.TransactionAt,
	}, nil
}
