package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret")

// Claims 定義 JWT 的負載
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Printf("[LOGIN ERROR] Binding JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "請檢查輸入資料"})
		return
	}

	// 模擬驗證
	if loginData.Username == "robby" && loginData.Password == "secret" {
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: loginData.Username,
			Role:     "admin", // 模擬 admin 角色
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Printf("[LOGIN ERROR] Token signing failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "伺服器錯誤"})
			return
		}

		log.Printf("[LOGIN SUCCESS] User %s logged in", loginData.Username)
		c.JSON(http.StatusOK, gin.H{
			"message": "登入成功",
			"token":   tokenString,
			"role":    claims.Role, // 將角色資訊返回
		})
		return
	}
	log.Printf("[LOGIN FAILED] Invalid credentials for user %s", loginData.Username)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
}

func LogoutHandler(c *gin.Context) {
	log.Println("[LOGOUT] User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func RefreshHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授權資訊"})
		return
	}
	var tokenString string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "授權資訊格式錯誤"})
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效或已過期的 token"})
		return
	}
	if time.Until(claims.ExpiresAt.Time) > 10*time.Minute {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token 尚未到刷新時間"})
		return
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 更新失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}
