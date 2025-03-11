package entities

import (
	"database/sql"
	"fmt"
)

type TaxesMigration struct {
    db *sql.DB
}

func NewTaxesMigration(db *sql.DB) *TaxesMigration {
    return &TaxesMigration{db: db}
}

func (m *TaxesMigration) Migrate() error {
    query := `
			CREATE TABLE IF NOT EXISTS tax_categories (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				description TEXT,
				status BOOLEAN, 
				is_default BOOLEAN,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP NULL
			);
    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Taxes table: %w", err)
    }
    
    fmt.Println("Taxes table ready")
    return nil
}