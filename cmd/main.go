package main

import (
	"log"
	"os"

	"go-book-learn/internal/db"
	"go-book-learn/internal/handlers"

	"github.com/gin-gonic/gin"
)

func initLogging() {
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	// 初始化日誌
	initLogging()

	// 初始化資料庫連線
	db.InitDB()

	// 建立 Gin 路由器
	r := gin.Default()

	// 載入所有路由
	handlers.RegisterRoutes(r)

	// 啟動伺服器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
