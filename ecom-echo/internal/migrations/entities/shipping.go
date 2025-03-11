package entities

import (
	"database/sql"
	"fmt"
)

type ShippingsMigration struct {
   db *sql.DB
}

func NewShippingsMigration(db *sql.DB) *ShippingsMigration {
   return &ShippingsMigration{db: db}
}

func (m *ShippingsMigration) Migrate() error {
   // Crear tabla shipping_methods
   createShippingMethodsQuery := `
       CREATE TABLE IF NOT EXISTS shipping_methods (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           description TEXT,
           code VARCHAR(50) NOT NULL UNIQUE,
           calculator_type VARCHAR(50),
           price_default DECIMAL(10,2),
           status BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_code (code),
           CONSTRAINT check_price CHECK (price_default >= 0)
       )
   `
   if _, err := m.db.Exec(createShippingMethodsQuery); err != nil {
       return fmt.Errorf("error creating shipping_methods table: %w", err)
   }

   // Crear tabla shipping_rules
   createShippingRulesQuery := `
       CREATE TABLE IF NOT EXISTS shipping_rules (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_shipping_method INT NOT NULL,
           max_value DECIMAL(10,2),
           min_value DECIMAL(10,2),
           price DECIMAL(10,2) NOT NULL,
           type VARCHAR(50) NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_shipping_method) REFERENCES shipping_methods(id),
           INDEX idx_shipping_method (id_shipping_method),
           CONSTRAINT check_rule_values CHECK (
               (max_value IS NULL OR min_value IS NULL) OR
               (max_value > min_value)
           ),
           CONSTRAINT check_rule_price CHECK (price >= 0)
       )
   `
   if _, err := m.db.Exec(createShippingRulesQuery); err != nil {
       return fmt.Errorf("error creating shipping_rules table: %w", err)
   }

   // Crear tabla eligibility_checkers
   createEligibilityCheckersQuery := `
       CREATE TABLE IF NOT EXISTS eligibility_checkers (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_shipping_method INT NOT NULL,
           type VARCHAR(50) NOT NULL,
           is_active BOOLEAN DEFAULT true,
           rule JSON,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_shipping_method) REFERENCES shipping_methods(id),
           INDEX idx_shipping_method (id_shipping_method),
           INDEX idx_type (type),
           INDEX idx_active (is_active)
       )
   `
   if _, err := m.db.Exec(createEligibilityCheckersQuery); err != nil {
       return fmt.Errorf("error creating eligibility_checkers table: %w", err)
   }

   // Crear tabla zones_shipping_methods
   createZonesShippingMethodsQuery := `
       CREATE TABLE IF NOT EXISTS n_zones_shipping_methods (
           id_shipping_method INT,
           id_zone INT,
           PRIMARY KEY (id_shipping_method, id_zone),
           FOREIGN KEY (id_shipping_method) REFERENCES shipping_methods(id),
           FOREIGN KEY (id_zone) REFERENCES zones(id),
           INDEX idx_zone (id_zone)
       )
   `
   if _, err := m.db.Exec(createZonesShippingMethodsQuery); err != nil {
       return fmt.Errorf("error creating zones_shipping_methods table: %w", err)
   }

   fmt.Println("Shipping system tables ready")
   return nil
}