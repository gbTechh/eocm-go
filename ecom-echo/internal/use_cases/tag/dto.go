package tag

import "time"

// CreateTagRequest estructura para crear un nuevo tag
type CreateTagRequest struct {
	Name   string `json:"name" validate:"required,min=2"`
	Code   string `json:"code" validate:"required,min=2"`
	Active bool   `json:"active"`
}

// TagResponse estructura para respuestas
type TagResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateTagRequest estructura para actualización
type UpdateTagRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2"`
	Code   *string `json:"code,omitempty" validate:"omitempty,min=2"`
	Active *bool   `json:"active,omitempty"`
}

// Pagination estructura para paginación y búsqueda
type Pagination struct {
	Page     int    `query:"page" validate:"min=1"`
	PerPage  int    `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
	Active   *bool  `query:"active"`
}