package category

import "time"

type Category struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Description string    `json:"description"`
    ParentID    *int64    `json:"parent_id"`
    IDMedia     *int64    `json:"id_media"`
    Media       *Media    `json:"media"`  // Para incluir la información completa del media
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Media es una versión simplificada de la entidad Media
type Media struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	MimeType  string    `json:"mime_type"`
	Size      float64   `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}