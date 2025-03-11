package media

import "time"

type Media struct {
	ID          int64     `json:"id"`
	FileName    string    `json:"file_name"`
	FilePath		string    `json:"file_path"`
	MimeType    string    `json:"mime_type"`
	Size        float64   `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}