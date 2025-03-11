package entities

import (
	"database/sql"
	"fmt"
)

type RolesMigration struct {
   db *sql.DB
}

func NewRolesMigration(db *sql.DB) *RolesMigration {
   return &RolesMigration{db: db}
}

func (m *RolesMigration) checkConstraintExists(constraintName string) (bool, error) {
   query := `
       SELECT COUNT(*)
       FROM information_schema.TABLE_CONSTRAINTS
       WHERE CONSTRAINT_SCHEMA = DATABASE()
       AND TABLE_NAME = 'users'
       AND CONSTRAINT_NAME = ?
   `
   var count int
   if err := m.db.QueryRow(query, constraintName).Scan(&count); err != nil {
       return false, fmt.Errorf("error checking constraint existence: %w", err)
   }
   return count > 0, nil
}

func (m *RolesMigration) Migrate() error {
   // Crear tabla roles
   createTableQuery := `
       CREATE TABLE IF NOT EXISTS roles (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(255) UNIQUE NOT NULL,
           is_root BOOLEAN DEFAULT FALSE,
           created_by INT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           deleted_at TIMESTAMP NULL,
           FOREIGN KEY (created_by) REFERENCES users(id)
       )
   `
   
   if _, err := m.db.Exec(createTableQuery); err != nil {
       return fmt.Errorf("error creating Roles table: %w", err)
   }

   // Verificar si la constraint existe
   exists, err := m.checkConstraintExists("fk_users_roles")
   if err != nil {
       return err
   }

   // Solo agregar la constraint si no existe
   if !exists {
       alterTableQuery := `
           ALTER TABLE users
           ADD CONSTRAINT fk_users_roles
           FOREIGN KEY (id_rol) REFERENCES roles(id)
       `
       
       if _, err := m.db.Exec(alterTableQuery); err != nil {
           return fmt.Errorf("error adding foreign key to Users table: %w", err)
       }
   }
   
   fmt.Println("Roles table and relationships ready")
   return nil
}