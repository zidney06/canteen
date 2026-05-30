package repository

import (
	"log"
	"server/src/config"
	"server/src/models"
	"server/src/types"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ClientRepo struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

func NewClientRepo(db *gorm.DB, rdb *redis.Client) *ClientRepo {
	return &ClientRepo{
		Db:  db,
		Rdb: rdb,
	}
}

func (r *ClientRepo) GetClientList() ([]types.Client, error) {
	var clients []types.Client

	result := r.Db.Model(&models.Client{}).Select("client_name", "id", "client_key").Find(&clients)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, &types.RepoError{
			Message: "Database error!",
			Type:    0,
		}
	}

	return clients, nil
}

func (r *ClientRepo) GetClientById(clientId string) (types.Client, error) {
	var client types.Client

	result := r.Db.Model(&models.Client{}).Select("id", "client_name", "client_key").First(&client, "id = ?", clientId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return types.Client{}, &types.RepoError{
				Message: "Client with " + clientId + " not found",
				Type:    1,
			}
		} else if result.Error != nil {
			log.Printf("%+v\n", result.Error)
			return types.Client{}, &types.RepoError{
				Message: "Database error!",
				Type:    0,
			}
		}
	}

	return client, nil
}

func (r *ClientRepo) GetClientSecretRepo(token string) (types.NewClientSecretType, error) {
	res, err := r.Rdb.HGetAll(config.Ctx, token).Result()
	if err != nil {
		log.Println(err.Error())
		return types.NewClientSecretType{}, &types.RepoError{
			Message: "Redis error!",
			Type:    0,
		}
	} else if len(res) == 0 {
		return types.NewClientSecretType{}, &types.RepoError{
			Message: "Record not found!",
			Type:    1,
		}
	}

	r.Rdb.Del(config.Ctx, token)

	return types.NewClientSecretType{
		ClientKey:    res["clientKey"],
		ClientSecret: res["clientSecret"],
	}, nil
}

func (r *ClientRepo) CreateNewClient(clientTamplate types.NewClientTamplate) (models.Client, error) {
	var newClient models.Client

	newClient.ID = clientTamplate.ID
	newClient.ClientName = clientTamplate.ClientName
	newClient.ClientSecret = clientTamplate.ClientSecret
	newClient.ClientKey = clientTamplate.ClientKey

	result := r.Db.Create(&newClient)
	if result.Error != nil {
		log.Println(result.Error)
		return models.Client{}, &types.RepoError{
			Message: "Database error",
			Type:    0,
		}
	}

	return newClient, nil
}

func (r *ClientRepo) SetToRedis(clientKey string, clientSecret string, token string) error {
	data := map[string]string{
		"clientKey":    clientKey,
		"clientSecret": clientSecret,
	}
	err := r.Rdb.HSet(config.Ctx, token, data).Err()
	if err != nil {
		log.Println(err)
		return &types.RepoError{
			Message: "Redis error!",
			Type:    0,
		}
	}

	// set expired karena HSet tidak bisa set exipred secara langsung
	expiration := 1 * time.Hour
	err = r.Rdb.Expire(config.Ctx, token, expiration).Err()
	if err != nil {
		log.Println(err)
		return &types.RepoError{
			Message: "Redis error!",
			Type:    0,
		}
	}

	return nil
}

func (r *ClientRepo) SetNewClientSecret(clientKey string, hashedNewSecret string) error {
	// rubah cs client di db
	result := r.Db.Model(&models.Client{}).Where("client_key = ?", clientKey).Update("client_secret", hashedNewSecret)

	// 4. Cek apakah ada error pada query (misal: masalah koneksi)
	if result.Error != nil {
		log.Println(result.Error)
		return &types.RepoError{
			Message: "Database error!",
			Type:    0,
		}
	} else if result.RowsAffected == 0 {
		return &types.RepoError{
			Message: "Client with " + clientKey + "not found!",
			Type:    1,
		}
	}

	return nil
}
