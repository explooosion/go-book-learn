package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "歡迎來到 Go 的 API 世界，這可是本小姐特別準備的喔！",
		})
	})

	// 登入路由
	r.POST("/login", loginHandler)

	r.Run(":8080")
}

// 登入處理函式
func loginHandler(c *gin.Context) {
	// 接收使用者傳入的 JSON 資料
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 綁定 JSON 到 loginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的資料有誤，重新檢查一下吧！"})
		return
	}

	// 模擬驗證帳號密碼
	if loginData.Username == "robby" && loginData.Password == "secret" {
		c.JSON(http.StatusOK, gin.H{"message": "登入成功～你應該感到榮幸，這可是本小姐批准的喔！"})

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤，你是故意惹本小姐生氣嗎？"})
	}
}
