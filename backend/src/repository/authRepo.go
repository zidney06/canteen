package repository

import (
	"errors"
	"fmt"
	"log"
	"server/src/config"
	"server/src/models"
	"server/src/types"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepo struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

func NewAuthRepo(db *gorm.DB, rdb *redis.Client) *AuthRepo {
	return &AuthRepo{
		Db:  db,
		Rdb: rdb,
	}
}

func (r *AuthRepo) GetClientWithKey(clientKey string) (types.ClientAuth, error) {
	var client types.ClientAuth

	result := r.Db.Model(&models.Client{}).Select("id", "client_secret", "client_key", "client_name").Where("client_key = ?", clientKey).First(&client)

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return types.ClientAuth{}, &types.RepoError{
				Message: "Client with key: " + clientKey + " not found!",
				Type:    2,
			}
		} else if result.Error != nil {
			log.Printf("%+v\n", result.Error)
			return types.ClientAuth{}, &types.RepoError{
				Message: "Database server error",
				Type:    0,
			}
		}
	}

	return client, nil
}

func (r *AuthRepo) IsTokenBlacklisted(claims types.MyClaims) (bool, string) {
	redisKey := fmt.Sprintf("blacklist:%s", claims.ID)

	result, err := r.Rdb.Exists(config.Ctx, redisKey).Result()
	if err != nil {
		log.Println(err)
		return false, "Something went wrong."
	}

	// result bernilai 1 jika ada, 0 jika tidak ada
	if result > 0 {
		return true, "Key already exist."
	} else {
		return false, "Key not found."
	}
}

func (r *AuthRepo) SetBlacklistedToken(claims types.MyClaims) error {
	expiration := 12 * time.Hour
	redisKey := fmt.Sprintf("blacklist:%s", claims.ID)

	err := r.Rdb.Set(config.Ctx, redisKey, "blacklisted", expiration).Err()
	if err != nil {
		log.Printf("%+v\n", err)
		return errors.New("Logout failed.")
	}

	return nil
}
