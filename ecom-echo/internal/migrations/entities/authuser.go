package entities

import (
	"database/sql"
	"fmt"
)

type AuthUserMigration struct {
   db *sql.DB
}

func NewAuthUserMigration(db *sql.DB) *AuthUserMigration {
   return &AuthUserMigration{db: db}
}

func (m *AuthUserMigration) Migrate() error {
   // Crear tabla LoginAttempt
   createLoginAttemptQuery := `
       CREATE TABLE IF NOT EXISTS login_attempts (
           id INT AUTO_INCREMENT PRIMARY KEY,
           email VARCHAR(255) NOT NULL,
           ip_address VARCHAR(45) NOT NULL,
           attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           success BOOLEAN DEFAULT FALSE,
           user_agent TEXT,
           failure_reason VARCHAR(255),
           INDEX idx_email (email),
           INDEX idx_ip_address (ip_address)
       )
   `
   
   if _, err := m.db.Exec(createLoginAttemptQuery); err != nil {
       return fmt.Errorf("error creating LoginAttempt table: %w", err)
   }

   // Crear tabla AccountLock
   createAccountLockQuery := `
       CREATE TABLE IF NOT EXISTS account_locks (
           id INT AUTO_INCREMENT PRIMARY KEY,
           email VARCHAR(255) NOT NULL,
           ip_address VARCHAR(45) NOT NULL,
           locked_until TIMESTAMP NOT NULL,
           reason VARCHAR(255),
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           INDEX idx_email (email),
           INDEX idx_ip_address (ip_address)
       )
   `
   
   if _, err := m.db.Exec(createAccountLockQuery); err != nil {
       return fmt.Errorf("error creating AccountLock table: %w", err)
   }

   // Crear tabla AuthSession
   createAuthSessionQuery := `
       CREATE TABLE IF NOT EXISTS auth_sessions (
           id INT AUTO_INCREMENT PRIMARY KEY,
           id_user INT NOT NULL,
           token VARCHAR(255) NOT NULL,
           ip_address VARCHAR(45),
           user_agent TEXT,
           last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           expires_at TIMESTAMP NOT NULL,
           is_active BOOLEAN DEFAULT TRUE,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           FOREIGN KEY (id_user) REFERENCES users(id),
           INDEX idx_token (token),
           INDEX idx_user (id_user)
       )
   `
   
   if _, err := m.db.Exec(createAuthSessionQuery); err != nil {
       return fmt.Errorf("error creating AuthSession table: %w", err)
   }
   
   fmt.Println("Auth tables ready")
   return nil
}