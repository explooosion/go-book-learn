package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "歡迎來到 Go 的 API 世界，這可是本小姐特別準備的喔！",
		})
	})

	r.Run(":8080")
}
