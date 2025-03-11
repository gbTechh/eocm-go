package currency

import "time"

type Currency struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Code      string    `json:"code"`
    Symbol    string    `json:"symbol"`
    IsBase    bool      `json:"is_base"`
    Active    bool      `json:"active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ExchangeRate struct {
    ID              int64     `json:"id"`
    FromCurrencyID  int64     `json:"from_currency_id"`
    ToCurrencyID    int64     `json:"to_currency_id"`
    Rate            float64   `json:"rate"`
    FromCurrency    *Currency `json:"from_currency,omitempty"`
    ToCurrency      *Currency `json:"to_currency,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
}