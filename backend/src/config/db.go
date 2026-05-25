package config

import (
	"log"
	"os"
	"server/src/models"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	onceDb sync.Once
)

func InitDB() {
	onceDb.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		isSqlite := os.Getenv("IS_SQLITE")
		var dbErr error

		if isSqlite == "1" {
			DB, dbErr = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		} else {
			DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}

		if dbErr != nil {
			log.Fatal("Gagal koneksi ke database:", dbErr)
		}

		if err := DB.AutoMigrate(&models.Client{}, &models.PurchaseItem{}, &models.Student{}, &models.Transaction{}); err != nil {
			log.Fatal("Gagal migrasi tabel:", err)
		}

		log.Println("Berhasil koneksi ke database!")
	})
}
