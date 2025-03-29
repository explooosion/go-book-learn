package handlers

import (
	"go-book-learn/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 認證路由
	r.POST("/login", LoginHandler)
	r.POST("/logout", LogoutHandler)
	r.POST("/refresh", RefreshHandler)

	// 產品路由（公開）
	r.GET("/products", GetProducts)
	r.GET("/products/:id", GetProductByID)

	// 受保護的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// 僅 admin 可操作的產品路由
		admin := protected.Group("/")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.POST("/products", CreateProduct)
			admin.PUT("/products/:id", UpdateProduct)
			admin.DELETE("/products/:id", DeleteProduct)
		}
	}
}
