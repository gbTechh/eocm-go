package entities

import (
	"database/sql"
	"fmt"
)

type ZonesMigration struct {
   db *sql.DB
}

func NewZonesMigration(db *sql.DB) *ZonesMigration {
   return &ZonesMigration{db: db}
}

func (m *ZonesMigration) Migrate() error {
   // Crear tabla zones
   createZonesQuery := `
       CREATE TABLE IF NOT EXISTS zones (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           active BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_active (active)
       )
   `
   if _, err := m.db.Exec(createZonesQuery); err != nil {
       return fmt.Errorf("error creating zones table: %w", err)
   }

   // Crear tabla tax_zones
   createTaxZonesQuery := `
       CREATE TABLE IF NOT EXISTS tax_zones (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_zone INT NOT NULL,
           name VARCHAR(100) NOT NULL,
           description TEXT,
           status BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_zone) REFERENCES zones(id),
           INDEX idx_zone (id_zone),
           INDEX idx_status (status)
       )
   `
   if _, err := m.db.Exec(createTaxZonesQuery); err != nil {
       return fmt.Errorf("error creating tax_zones table: %w", err)
   }

   // Crear tabla tax_rates
   createTaxRatesQuery := `
       CREATE TABLE IF NOT EXISTS tax_rates (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_tax_category INT NOT NULL,
           id_tax_zone INT NOT NULL,
           percentage DECIMAL(5,2) NOT NULL,
           is_default BOOLEAN DEFAULT false,
           status BOOLEAN DEFAULT true,
           name VARCHAR(100),
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_tax_category) REFERENCES tax_categories(id),
           FOREIGN KEY (id_tax_zone) REFERENCES tax_zones(id),
           INDEX idx_tax_category (id_tax_category),
           INDEX idx_tax_zone (id_tax_zone),
           CONSTRAINT check_percentage CHECK (percentage >= 0 AND percentage <= 100),
           CONSTRAINT unique_default_per_category UNIQUE (id_tax_category, is_default)
       )
   `
   if _, err := m.db.Exec(createTaxRatesQuery); err != nil {
       return fmt.Errorf("error creating tax_rates table: %w", err)
   }

   // Crear tabla tax_rules
   createTaxRulesQuery := `
       CREATE TABLE IF NOT EXISTS tax_rules (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_tax_rate INT NOT NULL,
           priority INT DEFAULT 0,
           status BOOLEAN DEFAULT true,
           min_amount DECIMAL(10,2),
           max_amount DECIMAL(10,2),
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_tax_rate) REFERENCES tax_rates(id),
           INDEX idx_tax_rate (id_tax_rate),
           INDEX idx_priority (priority),
           CONSTRAINT check_amounts CHECK (
               (max_amount IS NULL OR min_amount IS NULL) OR
               (max_amount > min_amount)
           ),
           CONSTRAINT check_min_amount CHECK (min_amount IS NULL OR min_amount >= 0),
           CONSTRAINT check_max_amount CHECK (max_amount IS NULL OR max_amount >= 0)
       )
   `
   if _, err := m.db.Exec(createTaxRulesQuery); err != nil {
       return fmt.Errorf("error creating tax_rules table: %w", err)
   }

   fmt.Println("Tax and zones system tables ready")
   return nil
}