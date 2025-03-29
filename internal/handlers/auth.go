package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// jwtKey 是用來簽署 JWT 的密鑰
var jwtKey = []byte("secret")

// Claims 定義 JWT 的負載
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"` // 使用者名稱
	Password string `json:"password"` // 密碼
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	Message string `json:"message"` // 成功訊息
	Token   string `json:"token"`   // JWT Token
	Role    string `json:"role"`    // 使用者角色
}

// RefreshResponse represents the refresh token response payload
type RefreshResponse struct {
	Token string `json:"token"` // 新生成的 JWT Token
}

// LogoutResponse represents the logout response payload
type LogoutResponse struct {
	Message string `json:"message"` // 登出成功訊息
}

// LoginHandler godoc
// @Summary 用戶登入
// @Description 驗證用戶的帳號與密碼，成功後返回 JWT token 與角色資訊
// @Tags 認證
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登入資訊"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var loginData LoginRequest
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
			"message": "登入成功～你真是太讓人心動了！",
			"token":   tokenString,
			"role":    claims.Role,
		})
		return
	}
	log.Printf("[LOGIN FAILED] Invalid credentials for user %s", loginData.Username)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤，你這樣可不行哦！"})
}

// LogoutHandler godoc
// @Summary 用戶登出
// @Description 執行登出操作，返回登出成功訊息
// @Tags 認證
// @Produce json
// @Success 200 {object} LogoutResponse
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	log.Println("[LOGOUT] User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "登出成功～期待再次相見哦！"})
}

// RefreshHandler godoc
// @Summary 刷新 JWT Token
// @Description 當 JWT Token 即將過期時，刷新生成新的 Token
// @Tags 認證
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /refresh [post]
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
