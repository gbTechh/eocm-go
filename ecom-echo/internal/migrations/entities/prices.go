package entities

import (
	"database/sql"
	"fmt"
)

type PricesMigration struct {
   db *sql.DB
}

func NewPricesMigration(db *sql.DB) *PricesMigration {
   return &PricesMigration{db: db}
}

func (m *PricesMigration) Migrate() error {
   // Crear tabla price_lists
   createPriceListsQuery := `
       CREATE TABLE IF NOT EXISTS price_lists (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           description TEXT,
           is_default BOOLEAN DEFAULT false,
           priority INT NOT NULL,
           status BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_priority (priority)           
       )
   `
   
   if _, err := m.db.Exec(createPriceListsQuery); err != nil {
       return fmt.Errorf("error creating PriceLists table: %w", err)
   }

   // Crear tabla prices
   createPricesQuery := `
       CREATE TABLE IF NOT EXISTS prices (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_price_list INT NOT NULL,
           amount DECIMAL(10,2) NOT NULL,
           starts_at TIMESTAMP NOT NULL,
           ends_at TIMESTAMP NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_price_list) REFERENCES price_lists(id),
           INDEX idx_price_list (id_price_list),
           INDEX idx_dates (starts_at, ends_at),
           CONSTRAINT check_amount CHECK (amount >= 0),
           CONSTRAINT check_dates CHECK (
               ends_at IS NULL OR starts_at < ends_at
           )
       )
   `
   
   if _, err := m.db.Exec(createPricesQuery); err != nil {
       return fmt.Errorf("error creating Prices table: %w", err)
   }

   // Crear tabla product_variant_prices
   createPriceVariantsQuery := `
       CREATE TABLE IF NOT EXISTS n_product_variant_prices (
           id_product_variant INT,
           id_price INT,
           is_active BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           PRIMARY KEY (id_price, id_product_variant),
           FOREIGN KEY (id_price) REFERENCES prices(id),
           FOREIGN KEY (id_product_variant) REFERENCES product_variants(id),
           INDEX idx_product_variant (id_product_variant),
           INDEX idx_active (is_active)
       )
   `
   
   if _, err := m.db.Exec(createPriceVariantsQuery); err != nil {
       return fmt.Errorf("error creating ProductVariant-Prices relationship table: %w", err)
   }
   
   fmt.Println("Price system tables ready")
   return nil
}