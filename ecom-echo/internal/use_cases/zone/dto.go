package zone

import "time"

// CreateZoneRequest estructura para crear un nuevo Zone
type CreateZoneRequest struct {
	Name   string `json:"name" validate:"required,min=2"`
	Active bool   `json:"active"`
}

// ZoneResponse estructura para respuestas
type ZoneResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateZoneRequest estructura para actualización
type UpdateZoneRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Active *bool   `json:"active,omitempty"`
}

// Pagination estructura para paginación y búsqueda
type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Active   *bool  `query:"active"`
}