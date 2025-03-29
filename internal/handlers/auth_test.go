package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestLoginHandlerSuccess 測試登入成功的情況
func TestLoginHandlerSuccess(t *testing.T) {
	router := gin.Default()
	router.POST("/login", LoginHandler)

	// 建立一個 LoginRequest 的 payload
	payload := LoginRequest{
		Username: "robby",
		Password: "secret",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 使用 httptest 建立一個響應記錄器
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 檢查狀態碼
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// 確認 token 與 role 存在
	assert.NotNil(t, resp["token"])
	assert.Equal(t, "admin", resp["role"])
}

// TestLoginHandlerFailure 測試登入失敗的情況
func TestLoginHandlerFailure(t *testing.T) {
	router := gin.Default()
	router.POST("/login", LoginHandler)

	payload := LoginRequest{
		Username: "robby",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
