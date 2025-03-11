// domain.go
package taxrule

import "time"

type TaxRule struct {
	ID          int64     `json:"id"`
	IDTaxRate   int64     `json:"id_tax_rate"`
	Priority    int       `json:"priority"`
	Status      bool      `json:"status"`
	MinAmount   *float64  `json:"min_amount"`
	MaxAmount   *float64  `json:"max_amount"`
	TaxRate     *TaxRate  `json:"tax_rate,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaxRate struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}
