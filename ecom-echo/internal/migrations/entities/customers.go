package entities

import (
	"database/sql"
	"fmt"
)

type CustomersMigration struct {
   db *sql.DB
}

func NewCustomersMigration(db *sql.DB) *CustomersMigration {
   return &CustomersMigration{db: db}
}

func (m *CustomersMigration) Migrate() error {
   // Crear tabla customers
   createCustomersQuery := `
       CREATE TABLE IF NOT EXISTS customers (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           last_name VARCHAR(100),
           email VARCHAR(100) UNIQUE NOT NULL,
           phone VARCHAR(20),
           active BOOLEAN DEFAULT true,
           is_guest BOOLEAN DEFAULT true,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           INDEX idx_email (email),
           INDEX idx_phone (phone)
       )
   `
   
   if _, err := m.db.Exec(createCustomersQuery); err != nil {
       return fmt.Errorf("error creating Customers table: %w", err)
   }

   // Crear tabla customer_groups
   createCustomerGroupsQuery := `
       CREATE TABLE IF NOT EXISTS customer_groups (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_price_list INT,
           name VARCHAR(100) NOT NULL,
           description VARCHAR(255),    
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_price_list) REFERENCES price_lists(id),
           INDEX idx_price (id_price_list)
       )
   `
   
   if _, err := m.db.Exec(createCustomerGroupsQuery); err != nil {
       return fmt.Errorf("error creating CustomerGroups table: %w", err)
   }

   // Crear tabla customers_groups
   createCustomersGroupsQuery := `
       CREATE TABLE IF NOT EXISTS n_customers_groups (
           id_customer INT,
           id_customer_group INT,
           PRIMARY KEY (id_customer, id_customer_group),
           FOREIGN KEY (id_customer) REFERENCES customers(id),
           FOREIGN KEY (id_customer_group) REFERENCES customer_groups(id),
           INDEX idx_customer_group (id_customer_group)
       )
   `
   
   if _, err := m.db.Exec(createCustomersGroupsQuery); err != nil {
       return fmt.Errorf("error creating N_Customer_Groups relationship table: %w", err)
   }

   // Crear tabla shipping_addresses
   createShippingAddressesQuery := `
       CREATE TABLE IF NOT EXISTS shipping_addresses (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_customer INT NOT NULL,
           id_city INT NOT NULL,
           full_name VARCHAR(100) NOT NULL,
           phone VARCHAR(20) NOT NULL,
           address_line1 VARCHAR(255) NOT NULL,
           address_line2 VARCHAR(255),
           postal_code VARCHAR(10),
           province VARCHAR(100),
           district VARCHAR(100),
           is_default BOOLEAN DEFAULT FALSE,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (id_customer) REFERENCES customers(id),
           FOREIGN KEY (id_city) REFERENCES cities(id),
           INDEX idx_customer (id_customer),
           INDEX idx_city (id_city),
           INDEX idx_postal_code (postal_code)
       )
   `
   
   if _, err := m.db.Exec(createShippingAddressesQuery); err != nil {
       return fmt.Errorf("error creating ShippingAddresses table: %w", err)
   }
   
   fmt.Println("Customers, groups and addresses ready")
   return nil
}