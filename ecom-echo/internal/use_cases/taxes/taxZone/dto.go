package taxzone

import "time"

type CreateTaxZoneRequest struct {
	IDZone      int64  `json:"id_zone" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

type TaxZoneResponse struct {
	ID          int64     `json:"id"`
	IDZone      int64     `json:"id_zone"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	Zone        *Zone     `json:"zone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateTaxZoneRequest struct {
	IDZone      *int64  `json:"id_zone,omitempty"`
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,name"`
	Description *string `json:"description,omitempty"`
	Status      *bool   `json:"status,omitempty"`
}

type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Status   *bool  `query:"status"`
}