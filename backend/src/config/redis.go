package config

import (
	"context"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var (
	Rdb     *redis.Client
	onceRdb sync.Once
)

func ConnectRdb() {
	onceRdb.Do(func() {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Username: "default",
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	})
}
