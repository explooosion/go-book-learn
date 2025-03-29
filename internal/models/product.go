package models

// Product 代表一個產品
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}
