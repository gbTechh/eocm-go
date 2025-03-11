package entities

import (
	"database/sql"
	"fmt"
)

type MediaMigration struct {
    db *sql.DB
}

func NewMediaMigration(db *sql.DB) *MediaMigration {
    return &MediaMigration{db: db}
}

func (m *MediaMigration) Migrate() error {
    query := `
        CREATE TABLE IF NOT EXISTS media (
            id INT AUTO_INCREMENT PRIMARY KEY,
            file_name VARCHAR(255) NOT NULL,
            file_path VARCHAR(255) NOT NULL,
            mime_type VARCHAR(20) NOT NULL,
            size DECIMAL(10, 6) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL
        );
    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Media table: %w", err)
    }
    
    fmt.Println("Media table ready")
    return nil
}