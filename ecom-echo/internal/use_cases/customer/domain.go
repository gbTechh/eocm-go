// domain.go
package customer

import (
	"time"
)

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Active    bool      `json:"active"`
	IsGuest   bool      `json:"is_guest"`
	Groups    []Group   `json:"groups,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Group struct {
	ID           int64      `json:"id"`
	IDPriceList  *int64     `json:"id_price_list,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	PriceList    *PriceList `json:"price_list,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PriceList struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsDefault   bool    `json:"is_default"`
	Priority    int     `json:"priority"`
}