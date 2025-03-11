package productattribute

import (
	attributeValue "ecom/internal/use_cases/attributeValue"
	"time"
)

// CreateProductAttribute estructura para crear un nuevo tag
type CreateProductAttributeRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=200,name"`
	Type     	  string `json:"type" validate:"required,min=2,max=100,name,oneof=select text number check radius"`
	Description string `json:"description" validate:"omitempty"`
}

// ProductAttributeResponse estructura para respuestas
type ProductAttributeResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
type ProductAttributeResponseList struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Values      []attributeValue.AttributeValue `json:"values"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UpdateProductAttribute estructura para actualización
type UpdateProductAttributeRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,name"`
	Type     	  *string `json:"type,omitempty" validate:"omitempty,name,oneof=select text number check radius"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
}


// CreateAttributeValue estructura para crear un nuevo attribute value
type CreateAttributeValueRequest struct {
	Name        				string `json:"name" validate:"required,min=2,max=200,name"`
	IDProductAttribute  int64  `json:"id_product_attribute" validate:"required"`
}

// AttributeValueResponse estructura para respuestas
type AttributeValueResponse struct {
	ID          							int64      `json:"id"`
	Name        							string     `json:"name"`
	IDProductAttribute        int64      `json:"id_product_attribute"`	
}

// UpdateAttributeValueRequest estructura para actualización
type UpdateAttributeValueRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,name"`
}

// Pagination estructura para paginación y búsqueda
type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
}