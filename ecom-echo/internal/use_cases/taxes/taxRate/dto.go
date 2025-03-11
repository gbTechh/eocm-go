// dto.go
package taxrate

import "time"

type CreateTaxRateRequest struct {
	IDTaxCategory int64   `json:"id_tax_category" validate:"required"`
	IDTaxZone     int64   `json:"id_tax_zone" validate:"required"`
	Percentage    float64 `json:"percentage" validate:"required,min=0,max=100"`
	IsDefault     bool    `json:"is_default"`
	Status        bool    `json:"status"`
	Name          string  `json:"name" validate:"required,min=2,name"`
}

type TaxRateResponse struct {
	ID            int64       `json:"id"`
	IDTaxCategory int64       `json:"id_tax_category"`
	IDTaxZone     int64       `json:"id_tax_zone"`
	Percentage    float64     `json:"percentage"`
	IsDefault     bool        `json:"is_default"`
	Status        bool        `json:"status"`
	Name          string      `json:"name"`
	TaxCategory   *TaxCategory `json:"tax_category,omitempty"`
	TaxZone      *TaxZone     `json:"tax_zone,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type UpdateTaxRateRequest struct {
	IDTaxCategory *int64   `json:"id_tax_category,omitempty"`
	IDTaxZone     *int64   `json:"id_tax_zone,omitempty"`
	Percentage    *float64 `json:"percentage,omitempty" validate:"omitempty,min=0,max=100"`
	IsDefault     *bool    `json:"is_default,omitempty"`
	Status        *bool    `json:"status,omitempty"`
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=2,name"`
}

type Pagination struct {
	Page     int     `query:"page" validate:"min=1"`
	PerPage  int     `query:"per_page" validate:"min=1,max=100"`
	Search   string  `query:"search"`
	Status   *bool   `query:"status"`
}