package entities

import (
	"database/sql"
	"fmt"
)

type TagsMigration struct {
    db *sql.DB
}

func NewTagsMigration(db *sql.DB) *TagsMigration {
    return &TagsMigration{db: db}
}

func (m *TagsMigration) Migrate() error {
    query := `
			CREATE TABLE IF NOT EXISTS tags (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				code VARCHAR(100) UNIQUE NOT NULL,
				active BOOLEAN, 
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP NULL
			);
    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Tags table: %w", err)
    }
    
    fmt.Println("Tags table ready")
    return nil
}