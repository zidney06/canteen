package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"server/src/models"
	"server/src/types"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StudentRepo struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

func NewStudentRepo(db *gorm.DB, rdb *redis.Client) *StudentRepo {
	return &StudentRepo{
		Db:  db,
		Rdb: rdb,
	}
}

func (r *StudentRepo) GetStudentById(studentId string) (types.GetStudent, error) {
	var std types.GetStudent

	result := r.Db.Model(&models.Student{}).Select("id", "name", "is_blocked").First(&std, "id = ?", studentId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return types.GetStudent{}, &types.RepoError{
				Message: "Client with " + studentId + " not found",
				Type:    1,
			}
		} else if result.Error != nil {
			log.Printf("%+v\n", result.Error)
			return types.GetStudent{}, &types.RepoError{
				Message: "Database error!",
				Type:    0,
			}
		}
	}

	return std, nil
}

func (r *StudentRepo) GetStudentFromApiById(studentId string) (types.StudentData, error) {
	baseUrl := os.Getenv("KELASKU_URL")
	url := baseUrl + "/" + studentId
	API_KEY := os.Getenv("KELASKU_API_KEY")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// ada 3 parameter, method, url, dan body. karena get gak pake body jadi isi nil aja
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request object: %v", err)
		return types.StudentData{}, &types.RepoError{
			Message: "Req failed!",
			Type:    0,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-key", API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
		return types.StudentData{}, &types.RepoError{
			Message: "Req failed!",
			Type:    0,
		}
	}
	defer resp.Body.Close() // Tutup body setelah digunakan

	// cek status request
	if resp.StatusCode == http.StatusNotFound {
		return types.StudentData{}, &types.RepoError{
			Message: "Data not found!",
			Type:    1,
		}
	} else if resp.StatusCode == http.StatusUnauthorized {
		return types.StudentData{}, &types.RepoError{
			Message: "Unauthorized!",
			Type:    4,
		}
	} else if resp.StatusCode != http.StatusOK {
		log.Print(resp)
		return types.StudentData{}, &types.RepoError{
			Message: "Error!",
			Type:    0,
		}
	}

	var result types.GetStudentByIdType
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return types.StudentData{}, &types.RepoError{
			Message: "Error while decode reeponse!",
			Type:    3,
		}
	}

	return result.Data, nil
}

func (r *StudentRepo) CreateNewStudent(studentId string) (types.Student, error) {
	if len(studentId) != 24 {
		return types.Student{}, &types.RepoError{
			Message: "Invalid student ID",
			Type:    1,
		}
	}

	// get ke api apk B untuk mendapatkan nama
	var repoErr *types.RepoError
	std, err := r.GetStudentFromApiById(studentId)
	if err != nil {
		if errors.As(err, &repoErr) {
			return types.Student{}, repoErr
		}
	}

	// simpan data ke db apk A
	var newStudentData models.Student

	newStudentData.ID = std.ID
	newStudentData.Name = std.Name

	result := r.Db.Create(&newStudentData)
	if result.Error != nil {
		log.Println(result.Error)
		return types.Student{}, &types.RepoError{
			Message: "Database error",
			Type:    0,
		}
	}

	// kiim hasil
	return types.Student{
		ID:        newStudentData.ID,
		Name:      newStudentData.Name,
		IsBlocked: newStudentData.IsBlocked,
	}, nil
}

// data sudah di enskripsi oleh servcie
func (r *StudentRepo) CreateQrCode(data string) (io.ReadCloser, error) {
	url := os.Getenv("QR_CODE_API")

	log.Println(url)

	fullURL := fmt.Sprintf("%s?text=%s&margin=2&format=svg", url, data)
	resp, err := http.Get(fullURL)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, &types.RepoError{
			Message: "Something went wrong!",
			Type:    0,
		}
	}
	return resp.Body, nil
}
