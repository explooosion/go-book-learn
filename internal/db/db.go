// internal/db/db.go
package db

import (
	"log"

	"go-book-learn/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化資料庫連線並進行自動遷移
func InitDB() {
	dsn := "root:secret@tcp(127.0.0.1:3306)/go_book_learn?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// 自動遷移 Product 資料結構對應的資料表
	DB.AutoMigrate(&models.Product{})
}
