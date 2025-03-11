package entities

import (
	"database/sql"
	"fmt"
)

type IndexesMigration struct {
    db *sql.DB
}

func NewIndexesMigration(db *sql.DB) *IndexesMigration {
    return &IndexesMigration{db: db}
}

func (m *IndexesMigration) Migrate() error {
    query := `
			CREATE INDEX IF NOT EXISTS idx_product_slug ON Products(slug);

    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Indexes constraint: %w", err)
    }
    
    fmt.Println("Indexes constraints ready")
    return nil
}