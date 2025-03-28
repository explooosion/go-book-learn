package main

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Product 代表一個簡單的產品資料結構
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

// 全域變數模擬資料庫
var products []Product
var nextID = 1

func main() {
	r := gin.Default()

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "歡迎來到 Go 的 API 世界，這可是本小姐特別準備的喔！",
		})
	})

	// 登入路由
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)

	// 產品 CRUD 路由
	r.GET("/products", getProducts)          // 取得所有產品
	r.GET("/products/:id", getProductByID)   // 依據 ID 取得單一產品
	r.POST("/products", createProduct)       // 新增產品
	r.PUT("/products/:id", updateProduct)    // 更新產品
	r.DELETE("/products/:id", deleteProduct) // 刪除產品

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

// 登出處理函式
func logoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "你已經成功登出囉～下次再來找本小姐吧！",
	})
}

// 取得所有產品
func getProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}

// 依據 ID 取得單一產品
func getProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}

// 新增產品
func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤，請檢查一下！"})
		return
	}
	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)
	c.JSON(http.StatusCreated, newProduct)
}

// 更新產品
func updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤，請檢查一下！"})
		return
	}

	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = p.ID // 保留原有的 ID
			products[i] = updatedProduct
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}

// 刪除產品
func deleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	for i, p := range products {
		if p.ID == id {
			products = slices.Delete(products, i, i+1)
			c.JSON(http.StatusOK, gin.H{"message": "產品已成功刪除～下次要小心一點喔！"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}
