// domain.go
package product

import (
	at "ecom/internal/use_cases/attributeValue"
	category "ecom/internal/use_cases/category"
	tag "ecom/internal/use_cases/tag"
	"time"
)

type Product struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Status      bool             `json:"status"`
	IDCategory  int64            `json:"id_category"`
	Category    *category.Category        `json:"category,omitempty"`
	Tags        []tag.Tag            `json:"tags,omitempty"`
	Variants    []ProductVariant `json:"variants,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ProductVariant struct {
	ID              int64            `json:"id"`
	Name            string           `json:"name"`
	SKU             string           `json:"sku"`
	Stock           int              `json:"stock"`
	IDProduct       int64            `json:"id_product"`
	IDTaxCategory   int64            `json:"id_tax_category"`
	AttributeValues []at.AttributeValue  `json:"attribute_values,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

