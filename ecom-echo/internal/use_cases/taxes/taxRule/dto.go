// dto.go
package taxrule

import "time"

type CreateTaxRuleRequest struct {
	IDTaxRate  int64    `json:"id_tax_rate" validate:"required"`
	Priority   int      `json:"priority"`
	Status     bool     `json:"status"`
	MinAmount  *float64 `json:"min_amount" validate:"omitempty,min=0"`
	MaxAmount  *float64 `json:"max_amount" validate:"omitempty,min=0,gtfield=MinAmount"`
}

type TaxRuleResponse struct {
	ID         int64     `json:"id"`
	IDTaxRate  int64     `json:"id_tax_rate"`
	Priority   int       `json:"priority"`
	Status     bool      `json:"status"`
	MinAmount  *float64  `json:"min_amount"`
	MaxAmount  *float64  `json:"max_amount"`
	TaxRate    *TaxRate  `json:"tax_rate,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateTaxRuleRequest struct {
	IDTaxRate  *int64   `json:"id_tax_rate,omitempty"`
	Priority   *int     `json:"priority,omitempty"`
	Status     *bool    `json:"status,omitempty"`
	MinAmount  *float64 `json:"min_amount,omitempty" validate:"omitempty,min=0"`
	MaxAmount  *float64 `json:"max_amount,omitempty" validate:"omitempty,min=0"`
}

type Pagination struct {
	Page       int    `query:"page" validate:"min=1"`
	PerPage    int    `query:"per_page" validate:"min=1,max=100"`
	Status     *bool  `query:"status"`
	IDTaxRate  *int64 `query:"id_tax_rate"`
}