package entities

import (
	"database/sql"
	"fmt"
)

type CurrenciesMigration struct {
   db *sql.DB
}

func NewCurrenciesMigration(db *sql.DB) *CurrenciesMigration {
   return &CurrenciesMigration{db: db}
}

func (m *CurrenciesMigration) Migrate() error {
   // Crear tabla currencies
   createCurrenciesQuery := `
       CREATE TABLE IF NOT EXISTS currencies (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(50) NOT NULL,
           code VARCHAR(3) NOT NULL UNIQUE,
           symbol VARCHAR(5),
           is_base BOOLEAN DEFAULT false,
           active BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_code (code)
       )
   `
   
   if _, err := m.db.Exec(createCurrenciesQuery); err != nil {
       return fmt.Errorf("error creating Currencies table: %w", err)
   }

   // Crear tabla actual_currencies
   createActualCurrencyQuery := `
        CREATE TABLE IF NOT EXISTS exchange_rates (
            id INT AUTO_INCREMENT PRIMARY KEY,
            from_currency_id INT,
            to_currency_id INT,
            rate DECIMAL(10,6) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (from_currency_id) REFERENCES currencies(id),
            FOREIGN KEY (to_currency_id) REFERENCES currencies(id),
            CONSTRAINT check_rate CHECK (rate > 0)
        );
   `
   
   if _, err := m.db.Exec(createActualCurrencyQuery); err != nil {
       return fmt.Errorf("error creating exchange_rates table: %w", err)
   }

   fmt.Println("Currencies, exchange rates and relationships ready")
   return nil
}