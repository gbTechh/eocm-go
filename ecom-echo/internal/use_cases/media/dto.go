package media

import "time"

// CreateMediaRequest es la estructura para subir un nuevo archivo
type CreateMediaRequest struct {
    File     []byte `json:"-" validate:"required"`         // Contenido del archivo
    FileName string `json:"file_name" validate:"required"` // Nombre original del archivo
    MimeType string `json:"mime_type" validate:"required,oneof=image/jpeg image/png image/gif application/pdf"`
}

// MediaResponse es la estructura de respuesta para información de archivos
type MediaResponse struct {
    ID        int64     `json:"id"`
    FileName  string    `json:"file_name"`
    FilePath  string    `json:"file_path"`
    MimeType  string    `json:"mime_type"`
    Size      float64   `json:"size"`
    URL       string    `json:"url"`        // URL pública del archivo
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UpdateMediaRequest para actualizar metadatos del archivo
type UpdateMediaRequest struct {
		File     []byte  `json:"-"`
    FilePath *string `json:"file_path,omitempty" validate:"omitempty,min=1"`
    FileName *string `json:"file_name,omitempty"`
    MimeType *string `json:"mime_type,omitempty"`
}

// MediaQuery para filtrar en búsquedas
type MediaQuery struct {
    MimeType []string `query:"mime_type,omitempty" validate:"omitempty,dive,oneof=image/jpeg image/png image/gif application/pdf"`
    FromDate string   `query:"from_date,omitempty" validate:"omitempty,datetime"`
    ToDate   string   `query:"to_date,omitempty" validate:"omitempty,datetime"`
    SizeMin  *float64 `query:"size_min,omitempty" validate:"omitempty,gte=0"`
    SizeMax  *float64 `query:"size_max,omitempty" validate:"omitempty,gte=0"`
    Page     int      `query:"page,omitempty" validate:"omitempty,gte=1"`
    PerPage  int      `query:"per_page,omitempty" validate:"omitempty,gte=1,lte=100"`
}
// MediaListResponse para paginación de resultados
type MediaListResponse struct {
    Items      []MediaResponse `json:"items"`
    TotalItems int64          `json:"total_items"`
    Page       int            `json:"page"`
    PerPage    int            `json:"per_page"`
    TotalPages int            `json:"total_pages"`
}

type Pagination struct {
	Page     int `query:"page" validate:"min=1"`
	PerPage  int `query:"per_page" validate:"min=1,max=100"`
	Search   string `query:"search"`
}