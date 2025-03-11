// internal/product/model.go
package product

import (
	"time"
)

type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CategoryID  string    `json:"category_id"`
    Category    *Category `json:"category,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type CreateProductInput struct {
    Name        string  `json:"name" binding:"required"`
    Description string  `json:"description" binding:"required"`
    Price       float64 `json:"price" binding:"required,gt=0"`
    Stock       int     `json:"stock" binding:"required,gte=0"`
    CategoryID  string  `json:"category_id" binding:"required"`
}

type UpdateProductInput struct {
    Name        *string  `json:"name,omitempty"`
    Description *string  `json:"description,omitempty"`
    Price       *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
    Stock       *int     `json:"stock,omitempty" binding:"omitempty,gte=0"`
    CategoryID  *string  `json:"category_id,omitempty"`
}