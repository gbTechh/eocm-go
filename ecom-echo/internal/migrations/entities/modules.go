package entities

import (
	"database/sql"
	"fmt"
)

type ModulesMigration struct {
    db *sql.DB
}

func NewModulesMigration(db *sql.DB) *ModulesMigration {
    return &ModulesMigration{db: db}
}

func (m *ModulesMigration) Migrate() error {
   // Crear tabla modules
   createModulesQuery := `
       CREATE TABLE IF NOT EXISTS modules (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(255) NOT NULL,
           code VARCHAR(255) UNIQUE NOT NULL,
           menu BOOLEAN DEFAULT FALSE,
           parent_id INT NULL,
           order_index INT NOT NULL,
           FOREIGN KEY (parent_id) REFERENCES modules(id)
       )
   `
   
   if _, err := m.db.Exec(createModulesQuery); err != nil {
       return fmt.Errorf("error creating Modules table: %w", err)
   }

   // Crear tabla rol_modules
   createRolModulesQuery := `
       CREATE TABLE IF NOT EXISTS rol_modules (
					id INT AUTO_INCREMENT PRIMARY KEY,
					role_id INT NOT NULL,
					module_id INT NOT NULL,
					can_view BOOLEAN DEFAULT FALSE,
					can_write BOOLEAN DEFAULT FALSE,
					can_edit BOOLEAN DEFAULT FALSE,
					can_delete BOOLEAN DEFAULT FALSE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					created_by INT NULL,
					FOREIGN KEY (role_id) REFERENCES roles(id),
					FOREIGN KEY (module_id) REFERENCES modules(id),
					FOREIGN KEY (created_by) REFERENCES users(id),
					UNIQUE KEY unique_role_module (role_id, module_id)
       )
   `
   
   if _, err := m.db.Exec(createRolModulesQuery); err != nil {
       return fmt.Errorf("error creating RolModules table: %w", err)
   }
   
   fmt.Println("Modules and RolModules tables ready")
   return nil
}