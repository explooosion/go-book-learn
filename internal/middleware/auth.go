package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret")

// Claims 定義與後端一致的 JWT 負載
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware 驗證 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供授權資訊"})
			return
		}
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "授權資訊格式錯誤"})
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無效或已過期的 token"})
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

// RoleMiddleware 檢查用戶角色
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := &Claims{}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供授權資訊"})
			return
		}
		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "授權資訊格式錯誤"})
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無效或已過期的 token"})
			return
		}
		if claims.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "權限不足"})
			return
		}
		c.Next()
	}
}
