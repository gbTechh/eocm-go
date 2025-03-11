package entities

import (
	"database/sql"
	"fmt"
)

type TriggersMigration struct {
    db *sql.DB
}

func NewTriggersMigration(db *sql.DB) *TriggersMigration {
    return &TriggersMigration{db: db}
}

func (m *TriggersMigration) Migrate() error {
    query := `
			CREATE TRIGGER IF NOT EXISTS delete_attribute_values
			AFTER UPDATE ON product_attributes
			FOR EACH ROW
			BEGIN
					IF NEW.deleted_at IS NOT NULL THEN
							DELETE FROM attribute_values 
							WHERE id_attribute_product = OLD.id;
					END IF;
			END;


    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Triggers constraint: %w", err)
    }
    
    fmt.Println("Triggers constraints ready")
    return nil
}