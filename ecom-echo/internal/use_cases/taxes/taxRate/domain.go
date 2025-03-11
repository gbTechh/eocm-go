// domain.go
package taxrate

import "time"

type TaxRate struct {
    ID             int64      `json:"id"`
    IDTaxCategory  int64      `json:"id_tax_category"`
    IDTaxZone      int64      `json:"id_tax_zone"`
    Percentage     float64    `json:"percentage"`
    IsDefault      bool       `json:"is_default"`
    Status         bool       `json:"status"`
    Name           string     `json:"name"`
    TaxCategory    *TaxCategory `json:"tax_category,omitempty"`
    TaxZone       *TaxZone    `json:"tax_zone,omitempty"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}

type TaxCategory struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type TaxZone struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}