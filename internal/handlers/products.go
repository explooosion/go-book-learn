package handlers

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"go-book-learn/internal/models"

	"github.com/gin-gonic/gin"
)

// 全域變數模擬資料庫，並加入 Mutex 保護
var (
	products []models.Product
	nextID   int
	mu       sync.RWMutex // 使用讀寫鎖保護全域變數
)

// GetProducts godoc
// @Summary 取得所有產品
// @Description 返回所有產品列表，這是公開 API，若沒有產品則返回空陣列
// @Tags 產品
// @Produce json
// @Success 200 {array} models.Product
// @Router /products [get]
func GetProducts(c *gin.Context) {
	log.Println("[GET PRODUCTS] Fetching products")
	mu.RLock()
	defer mu.RUnlock()
	// 如果 products 為 nil，則返回空的切片，避免 JSON 序列化為 null
	var prods []models.Product
	if products == nil {
		prods = []models.Product{}
	} else {
		prods = products
	}
	c.JSON(http.StatusOK, prods)
}

// GetProductByID godoc
// @Summary 取得單一產品
// @Description 根據產品 ID 返回單一產品資料
// @Tags 產品
// @Produce json
// @Param id path int true "產品 ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[GET PRODUCT ERROR] Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	mu.RLock()
	defer mu.RUnlock()
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}

// CreateProduct godoc
// @Summary 新增產品
// @Description 新增一個產品。此操作需要 admin 權限（JWT Token）
// @Tags 產品
// @Accept json
// @Produce json
// @Param product body models.Product true "產品資訊"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		log.Printf("[CREATE PRODUCT ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤"})
		return
	}

	mu.Lock()
	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)
	mu.Unlock()

	log.Printf("[CREATE PRODUCT] Product ID %d created", newProduct.ID)
	c.JSON(http.StatusCreated, newProduct)
}

// UpdateProduct godoc
// @Summary 更新產品
// @Description 根據產品 ID 更新產品資料，需提供完整資料（需要 admin 權限）
// @Tags 產品
// @Accept json
// @Produce json
// @Param id path int true "產品 ID"
// @Param product body models.Product true "更新後的產品資料"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = p.ID // 保留原有的 ID
			products[i] = updatedProduct
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}

// DeleteProduct godoc
// @Summary 刪除產品
// @Description 根據產品 ID 刪除產品資料（需要 admin 權限）
// @Tags 產品
// @Produce json
// @Param id path int true "產品 ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, p := range products {
		if p.ID == id {
			products = slices.Delete(products, i, i+1)
			c.JSON(http.StatusOK, gin.H{"message": "產品已刪除"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}
