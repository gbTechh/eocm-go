package entities

import (
	"database/sql"
	"fmt"
)

type CartsMigration struct {
   db *sql.DB
}

func NewCartsMigration(db *sql.DB) *CartsMigration {
   return &CartsMigration{db: db}
}

func (m *CartsMigration) Migrate() error {
   // Crear tabla carts
   createCartsQuery := `
       CREATE TABLE IF NOT EXISTS carts (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_customer INT NOT NULL,
           status VARCHAR(50) NOT NULL,
           session_id VARCHAR(255),
           total_amount DECIMAL(10,2) DEFAULT 0,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_customer) REFERENCES customers(id),
           INDEX idx_customer (id_customer),
           INDEX idx_session (session_id)
       )
   `
   
   if _, err := m.db.Exec(createCartsQuery); err != nil {
       return fmt.Errorf("error creating Carts table: %w", err)
   }

   // Crear tabla cart_items
   createCartItemsQuery := `
       CREATE TABLE IF NOT EXISTS cart_items (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_cart INT NOT NULL,
           id_product_variant INT NOT NULL,
           quantity INT NOT NULL,
           unit_price DECIMAL(10,2) NOT NULL,
           total_price DECIMAL(10,2) NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           FOREIGN KEY (id_cart) REFERENCES carts(id),
           FOREIGN KEY (id_product_variant) REFERENCES product_variants(id),
           INDEX idx_cart (id_cart),
           INDEX idx_product_variant (id_product_variant),
           CONSTRAINT check_quantity CHECK (quantity > 0),
           CONSTRAINT check_prices CHECK (unit_price >= 0 AND total_price >= 0)
       )
   `
   
   if _, err := m.db.Exec(createCartItemsQuery); err != nil {
       return fmt.Errorf("error creating CartItems table: %w", err)
   }
   
   fmt.Println("Carts and CartItems tables ready")
   return nil
}