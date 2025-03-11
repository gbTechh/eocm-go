// domain.go
package prices

import (
	"time"
)

type PriceList struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsDefault   bool      `json:"is_default"`
	Priority    int       `json:"priority"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Price struct {
	ID           int64      `json:"id"`
	IDPriceList  int64      `json:"id_price_list"`
	Amount       float64    `json:"amount"`
	StartsAt     time.Time  `json:"starts_at"`
	EndsAt       *time.Time `json:"ends_at,omitempty"`
	PriceList    *PriceList `json:"price_list,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ProductVariantPrice struct {
	IDProductVariant int64     `json:"id_product_variant"`
	IDPrice         int64     `json:"id_price"`
	IsActive        bool      `json:"is_active"`
	Price           *Price    `json:"price,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
