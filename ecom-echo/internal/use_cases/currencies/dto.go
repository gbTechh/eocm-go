package currency

import "time"

type CreateCurrencyRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Code    string `json:"code" validate:"required,len=3"`
	Symbol  string `json:"symbol" validate:"required"`
	IsBase  bool   `json:"is_base"`
	Active  bool   `json:"active"`
}

type UpdateCurrencyRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Code    *string `json:"code,omitempty" validate:"omitempty,len=3"`
	Symbol  *string `json:"symbol,omitempty"`
	IsBase  *bool   `json:"is_base,omitempty"`
	Active  *bool   `json:"active,omitempty"`
}

type CreateExchangeRateRequest struct {
	FromCurrencyID int64   `json:"from_currency_id" validate:"required"`
	ToCurrencyID   int64   `json:"to_currency_id" validate:"required"`
	Rate           float64 `json:"rate" validate:"required,gt=0"`
}

type CurrencyResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Symbol    string    `json:"symbol"`
	IsBase    bool      `json:"is_base"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExchangeRateResponse struct {
	ID           int64     `json:"id"`
	FromCurrency *Currency `json:"from_currency"`
	ToCurrency   *Currency `json:"to_currency"`
	Rate         float64   `json:"rate"`
	CreatedAt    time.Time `json:"created_at"`
}

type Pagination struct {
	Page    int    `query:"page" validate:"min=1"`
	PerPage int    `query:"per_page" validate:"min=1,max=100"`
	Search  string `query:"search"`
	Active  *bool  `query:"active"`
}

type ConvertAmountRequest struct {
	Amount            float64 `json:"amount" validate:"required,gt=0"`
	FromCurrencyCode string  `json:"from_currency_code" validate:"required,len=3"`
	ToCurrencyCode   string  `json:"to_currency_code" validate:"required,len=3"`
}

type ConvertAmountResponse struct {
	OriginalAmount   float64  `json:"original_amount"`
	ConvertedAmount  float64  `json:"converted_amount"`
	FromCurrency     Currency `json:"from_currency"`
	ToCurrency       Currency `json:"to_currency"`
	ExchangeRate     float64  `json:"exchange_rate"`
	ConvertedAt      time.Time `json:"converted_at"`
}