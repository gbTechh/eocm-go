package entities

import (
	"database/sql"
	"fmt"
)

type CategoriesMigration struct {
    db *sql.DB
}

func NewCategoriesMigration(db *sql.DB) *CategoriesMigration {
    return &CategoriesMigration{db: db}
}

func (m *CategoriesMigration) Migrate() error {
    query := `
        CREATE TABLE IF NOT EXISTS categories (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(200) NOT NULL,
            slug VARCHAR(220) UNIQUE NOT NULL,
            description TEXT,
            parent_id INT NULL,  -- Para subcategorías
            id_media INT NULL,   -- Imagen principal
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
            deleted_at TIMESTAMP NULL,
            FOREIGN KEY (parent_id) REFERENCES categories(id),
            FOREIGN KEY (id_media) REFERENCES media(id)
        );
    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating categories table: %w", err)
    }
    
    fmt.Println("Categories table ready")
    return nil
}