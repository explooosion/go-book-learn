package main

import (
	"log"
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
		log.Printf("[LOGIN ERROR] Binding JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的資料有誤，重新檢查一下吧！"})
		return
	}

	// 模擬驗證帳號密碼
	if loginData.Username == "robby" && loginData.Password == "secret" {
		log.Printf("[LOGIN SUCCESS] User %s logged in successfully", loginData.Username)
		c.JSON(http.StatusOK, gin.H{"message": "登入成功～你應該感到榮幸，這可是本小姐批准的喔！"})

	} else {
		log.Printf("[LOGIN FAILED] Invalid credentials for user %s", loginData.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤，你是故意惹本小姐生氣嗎？"})
	}
}

// 登出處理函式
func logoutHandler(c *gin.Context) {
	log.Println("[LOGOUT] User logged out")
	c.JSON(http.StatusOK, gin.H{
		"message": "你已經成功登出囉～下次再來找本小姐吧！",
	})
}

// 取得所有產品
func getProducts(c *gin.Context) {
	log.Println("[GET PRODUCTS] Fetching all products")
	c.JSON(http.StatusOK, products)
}

// 依據 ID 取得單一產品
func getProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[GET PRODUCT ERROR] Invalid product ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}
	for _, p := range products {
		if p.ID == id {
			log.Printf("[GET PRODUCT] Found product with ID %d", id)
			c.JSON(http.StatusOK, p)
			return
		}
	}
	log.Printf("[GET PRODUCT ERROR] Product with ID %d not found", id)
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}

// 新增產品
func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		log.Printf("[CREATE PRODUCT ERROR] Binding JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤，請檢查一下！"})
		return
	}
	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)
	log.Printf("[CREATE PRODUCT] New product added with ID %d", newProduct.ID)
	c.JSON(http.StatusCreated, newProduct)
}

// 更新產品
func updateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[UPDATE PRODUCT ERROR] Invalid product ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	var updatedProduct Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		log.Printf("[UPDATE PRODUCT ERROR] Binding JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤，請檢查一下！"})
		return
	}

	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = p.ID // 保留原有的 ID
			products[i] = updatedProduct
			log.Printf("[UPDATE PRODUCT] Product with ID %d updated", id)
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}
	log.Printf("[UPDATE PRODUCT ERROR] Product with ID %d not found", id)
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}

// 刪除產品
func deleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[DELETE PRODUCT ERROR] Invalid product ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	for i, p := range products {
		if p.ID == id {
			// 使用 slices.Delete (需 Go 1.21+)
			products = slices.Delete(products, i, i+1)
			log.Printf("[DELETE PRODUCT] Product with ID %d deleted", id)
			c.JSON(http.StatusOK, gin.H{"message": "產品已成功刪除～下次要小心一點喔！"})
			return
		}
	}
	log.Printf("[DELETE PRODUCT ERROR] Product with ID %d not found", id)
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品，難道是本小姐不小心賣掉了嗎？"})
}
