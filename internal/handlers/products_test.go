package handlers

import (
	"bytes"
	"encoding/json"
	"go-book-learn/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// generateAdminToken 用來產生一個有效的 admin JWT token
func generateAdminToken() string {
	expirationTime := time.Now().Add(5 * time.Minute)
	// 注意：此處的 Claims 結構需與你後端 auth.go 中的一致
	claims := &Claims{
		Username: "robby",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return "Bearer " + tokenString
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router)
	return router
}

// TestGetProducts 測試取得產品列表
func TestGetProducts(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var products []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &products)
	assert.NoError(t, err)
	// 初始情況可能為空
}

// TestCreateProduct 測試新增產品（模擬 admin token）
// 為了測試這個功能，我們可以直接模擬一個有效的 token，但這裡只示範請求流程
func TestCreateProduct(t *testing.T) {
	router := setupRouter() // 你的 setupRouter() 用來載入所有路由
	adminToken := generateAdminToken()

	newProd := models.Product{
		Name:  "Test Product",
		Price: 99.99,
	}
	body, _ := json.Marshal(newProd)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", adminToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var prodResp models.Product
	err := json.Unmarshal(w.Body.Bytes(), &prodResp)
	assert.NoError(t, err)
	assert.Equal(t, newProd.Name, prodResp.Name)
	assert.Equal(t, newProd.Price, prodResp.Price)
}

// 其他測試，如 UpdateProduct、DeleteProduct，也可類似撰寫
