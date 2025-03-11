// dto.go
package prices

import "time"

type CreatePriceListRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	Priority    int    `json:"priority" validate:"required,min=0"`
	Status      bool   `json:"status"`
}

type UpdatePriceListRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
	Priority    *int    `json:"priority,omitempty" validate:"omitempty,min=0"`
	Status      *bool   `json:"status,omitempty"`
}

type CreatePriceRequest struct {
	IDPriceList int64      `json:"id_price_list" validate:"required"`
	Amount      float64    `json:"amount" validate:"required,min=0"`
	StartsAt    time.Time  `json:"starts_at" validate:"required"`
	EndsAt      *time.Time `json:"ends_at,omitempty" validate:"omitempty,gtfield=StartsAt"`
}

type UpdatePriceRequest struct {
	Amount   *float64    `json:"amount,omitempty" validate:"omitempty,min=0"`
	StartsAt *time.Time  `json:"starts_at,omitempty"`
	EndsAt   *time.Time  `json:"ends_at,omitempty" validate:"omitempty,gtfield=StartsAt"`
}

type AssignPriceRequest struct {
	IDProductVariant int64 `json:"id_product_variant" validate:"required"`
	IDPrice         int64 `json:"id_price" validate:"required"`
	IsActive        bool  `json:"is_active"`
}

type PriceListResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"is_default"`
	Priority    int       `json:"priority"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PriceResponse struct {
	ID          int64            `json:"id"`
	IDPriceList int64            `json:"id_price_list"`
	Amount      float64          `json:"amount"`
	StartsAt    time.Time        `json:"starts_at"`
	EndsAt      *time.Time       `json:"ends_at,omitempty"`
	PriceList   *PriceListResponse `json:"price_list,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Status   *bool  `query:"status"`
}