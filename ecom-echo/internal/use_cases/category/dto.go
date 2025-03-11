package category

import "time"

// CreateTagRequest estructura para crear un nuevo tag
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=200,name"`
	Slug     	  string `json:"slug" validate:"required,min=2,max=200,slug"`
	ParentID    *int64  `json:"parent_id" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	IDMedia     *int64  `json:"id_media"`
}

// TagResponse estructura para respuestas
type CategoryResponse struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	ParentID    *int64         `json:"parent_id"`
	IDMedia     *int64         `json:"id_media"`
	Media       *Media         `json:"media,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UpdateTagRequest estructura para actualización
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,omitempty,min=2,max=200,name"`
	Slug     	  *string `json:"slug,omitempty" validate:"omitempty,omitempty,min=2,max=200,slug"`
	ParentID    *int64  `json:"parent_id,omitempty" validate:"omitempty,omitempty"`
	Description *string `json:"description,omitempty" validate:"omitempty"`
	IDMedia     *int64  `json:"id_media,omitempty"`
}

// Pagination estructura para paginación y búsqueda
type Pagination struct {
    Page     int    `query:"page" validate:"min=1"`
    PerPage  int    `query:"per_page" validate:"min=1,max=100"`
    Search   string `query:"search"`
    ParentID *int64 `query:"parent_id"`
}