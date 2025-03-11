// dto.go
package product

import (
	category "ecom/internal/use_cases/category"
	tag "ecom/internal/use_cases/tag"
	"time"
)

type CreateProductRequest struct {
	Name              string  `json:"name" validate:"required,min=2"`
	Slug              string  `json:"slug" validate:"required,min=2,slug"`
	Description       string  `json:"description"`
	Status           bool    `json:"status"`
	IDCategory       int64   `json:"id_category" validate:"required"`
	TagIDs           []int64 `json:"tag_ids"`
	DefaultVariant   CreateProductVariantRequest `json:"default_variant" validate:"required"`
}

type CreateProductVariantRequest struct {
	Name            string   `json:"name" validate:"omitempty,name"`
	SKU             string   `json:"sku" validate:"omitempty,name"`
	Stock           int      `json:"stock" validate:"min=0"`
	IDTaxCategory   int64    `json:"id_tax_category" validate:"omitempty"`
	AttributeValues []CreateAttributeValueRequest `json:"attribute_values"`
}

type CreateAttributeValueRequest struct {
	ID                 int64  `json:"id" validate:"required"`
}

type ProductResponse struct {
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

type UpdateProductRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,min=2,slug"`
	Description *string `json:"description,omitempty"`
	Status      *bool   `json:"status,omitempty"`
	IDCategory  *int64  `json:"id_category,omitempty"`
}

type CreateVariantRequest struct {
	Name            string   `json:"name" validate:"name"`
	SKU             string   `json:"sku"`
	Stock           int      `json:"stock" validate:"min=0"`
	IDTaxCategory   int64    `json:"id_tax_category" validate:"omitempty"`
	IDProduct       int64    `json:"id_product" validate:"required"`
	AttributeValues []CreateAttributeValueRequest `json:"attribute_values"`
}

type UpdateVariantRequest struct {
	Name          *string `json:"name,omitempty"`
	SKU           *string `json:"sku,omitempty"`
	Stock         *int    `json:"stock,omitempty" validate:"omitempty,min=0"`
	IDTaxCategory *int64  `json:"id_tax_category,omitempty"`
	AttributeValues *[]CreateAttributeValueRequest `json:"attribute_values,omitempty"`
	
}

type Pagination struct {
	Page       int    `query:"page" validate:"min=1"`
	PerPage    int    `query:"per_page" validate:"min=1,max=100"`
	Search     string `query:"search"`
	Status     *bool  `query:"status"`
	IDCategory *int64 `query:"id_category"`
}

