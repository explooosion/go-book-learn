// @title Go Book Learn API
// @version 1.0
// @description This is a sample server for Go Book Learn.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
package main

import (
	_ "go-book-learn/docs"
	"go-book-learn/internal/db"
	"go-book-learn/internal/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func initLogging() {
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	initLogging()
	db.InitDB()
	r := gin.Default()

	// Register all routes including Swagger docs later
	handlers.RegisterRoutes(r)

	// Swagger docs route
	// 注意：這裡使用 gin-swagger 與 swaggo/files
	// @Router /swagger/*any [get]
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
