package handlers

import (
	"log"
	"net/http"
	"slices"
	"strconv"

	"go-book-learn/internal/models"

	"github.com/gin-gonic/gin"
)

var products = []models.Product{}
var nextID = 1

func GetProducts(c *gin.Context) {
	log.Println("[GET PRODUCTS] Fetching products")
	c.JSON(http.StatusOK, products)
}

func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[GET PRODUCT ERROR] Invalid ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}

func CreateProduct(c *gin.Context) {
	var newProduct models.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		log.Printf("[CREATE PRODUCT ERROR] %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "輸入的產品資料有誤"})
		return
	}
	newProduct.ID = nextID
	nextID++
	products = append(products, newProduct)
	log.Printf("[CREATE PRODUCT] Product ID %d created", newProduct.ID)
	c.JSON(http.StatusCreated, newProduct)
}

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
	for i, p := range products {
		if p.ID == id {
			updatedProduct.ID = p.ID
			products[i] = updatedProduct
			c.JSON(http.StatusOK, updatedProduct)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}

func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的產品 ID"})
		return
	}
	for i, p := range products {
		if p.ID == id {
			products = slices.Delete(products, i, i+1)
			c.JSON(http.StatusOK, gin.H{"message": "產品已刪除"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "找不到該產品"})
}
