package service

import (
	"io"
	"log"
	"os"
	"server/src/repository"
	"server/src/types"
	"server/src/utils"
)

type StudentService struct {
	Repo *repository.StudentRepo
}

func NewStudentService(repo *repository.StudentRepo) *StudentService {
	return &StudentService{
		Repo: repo,
	}
}

func (s *StudentService) CreateQrCode(studentId string) (io.ReadCloser, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY")) // 32 bytes

	// Proses Enkripsi
	encryptedText, err := utils.Encrypt(studentId, secretKey)
	if err != nil {
		return nil, &types.ServiceError{
			Message:    "Failed to encrypt student id",
			HttpStatus: 500,
		}
	}
	log.Println("encrypted: ", encryptedText)

	// Panggil repo, terima stream-nya, lalu oper ke atas
	stream, err := s.Repo.CreateQrCode(encryptedText)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
