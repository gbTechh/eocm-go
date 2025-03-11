package entities

import (
	"database/sql"
	"fmt"
)

type OrdersMigration struct {
   db *sql.DB
}

func NewOrdersMigration(db *sql.DB) *OrdersMigration {
   return &OrdersMigration{db: db}
}

func (m *OrdersMigration) Migrate() error {
   // Crear tabla payment_methods 
   createPaymentMethodsQuery := `
       CREATE TABLE IF NOT EXISTS payment_methods (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           code VARCHAR(50) NOT NULL UNIQUE,
           description TEXT,
           status BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_code (code)
       )
   `
   
   if _, err := m.db.Exec(createPaymentMethodsQuery); err != nil {
       return fmt.Errorf("error creating PaymentMethods table: %w", err)
   }

   // Crear tabla orders
   createOrdersQuery := `
       CREATE TABLE IF NOT EXISTS orders (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_customer INT NOT NULL,
           id_shipping_address INT NOT NULL,
           id_shipping_method INT NOT NULL,
           id_payment_method INT NOT NULL,
           id_currency INT NOT NULL,
           order_number VARCHAR(50) UNIQUE NOT NULL,
           status VARCHAR(50) NOT NULL,
           subtotal DECIMAL(10,2) NOT NULL,
           shipping_cost DECIMAL(10,2) NOT NULL,
           tax_amount DECIMAL(10,2) NOT NULL,
           total_amount DECIMAL(10,2) NOT NULL,
           total_items INT NOT NULL,
           tracking_number VARCHAR(100),
           notes TEXT,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           cancelled_at TIMESTAMP NULL,
           cancelled_reason VARCHAR(255),
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_customer) REFERENCES customers(id),
           FOREIGN KEY (id_shipping_address) REFERENCES shipping_addresses(id),
           FOREIGN KEY (id_shipping_method) REFERENCES shipping_methods(id),
           FOREIGN KEY (id_payment_method) REFERENCES payment_methods(id),
           FOREIGN KEY (id_currency) REFERENCES currencies(id),
           INDEX idx_customer (id_customer),
           INDEX idx_order_number (order_number),
           INDEX idx_status (status),
           CONSTRAINT check_amounts_orders CHECK (
               subtotal >= 0 AND 
               shipping_cost >= 0 AND 
               tax_amount >= 0 AND 
               total_amount >= 0 AND
               total_items >= 0
           )
       )
   `
   
   if _, err := m.db.Exec(createOrdersQuery); err != nil {
       return fmt.Errorf("error creating Orders table: %w", err)
   }

   // Crear tabla order_items
   createOrderItemsQuery := `
       CREATE TABLE IF NOT EXISTS order_items (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_order INT NOT NULL,
           id_product_variant INT NOT NULL,
           quantity INT NOT NULL,
           unit_price DECIMAL(10,2) NOT NULL,
           tax_amount DECIMAL(10,2) NOT NULL,
           total_price DECIMAL(10,2) NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           FOREIGN KEY (id_order) REFERENCES orders(id),
           FOREIGN KEY (id_product_variant) REFERENCES product_variants(id),
           INDEX idx_order (id_order),
           INDEX idx_product_variant (id_product_variant),
           CONSTRAINT check_item_quantity CHECK (quantity > 0),
           CONSTRAINT check_item_prices CHECK (
               unit_price >= 0 AND 
               tax_amount >= 0 AND 
               total_price >= 0
           )
       )
   `
   
   if _, err := m.db.Exec(createOrderItemsQuery); err != nil {
       return fmt.Errorf("error creating OrderItems table: %w", err)
   }

   // Crear tabla payments
   createPaymentsQuery := `
       CREATE TABLE IF NOT EXISTS payments (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_order INT NOT NULL,
           id_payment_method INT NOT NULL,
           amount DECIMAL(10,2) NOT NULL,
           status VARCHAR(50) NOT NULL,
           transaction_id VARCHAR(100),
           external_reference VARCHAR(255),
           payment_data JSON,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_order) REFERENCES orders(id),
           FOREIGN KEY (id_payment_method) REFERENCES payment_methods(id),
           INDEX idx_order (id_order),
           INDEX idx_transaction (transaction_id),
           INDEX idx_status (status),
           CONSTRAINT check_payment_amount CHECK (amount > 0)
       )
   `
   
   if _, err := m.db.Exec(createPaymentsQuery); err != nil {
       return fmt.Errorf("error creating Payments table: %w", err)
   }
   
   fmt.Println("Orders system tables ready")
   return nil
}