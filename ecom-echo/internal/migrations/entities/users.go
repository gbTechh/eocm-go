package entities

import (
	"database/sql"
	"fmt"
)

type UsersMigration struct {
    db *sql.DB
}

func NewUsersMigration(db *sql.DB) *UsersMigration {
    return &UsersMigration{db: db}
}

func (m *UsersMigration) Migrate() error {
    query := `
			CREATE TABLE IF NOT EXISTS users (
				id INt AUTO_INCREMENT PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				id_rol INT NOT NULL,
				password VARCHAR(255) NOT NULL,
				name VARCHAR(255) NOT NULL,
				state BOOLEAN DEFAULT FALSE,
				last_login TIMESTAMP NULL,
				failed_attempts INT DEFAULT 0,
				last_failed_attempt TIMESTAMP NULL,
				password_changed_at TIMESTAMP NULL,
				force_password_change BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP NULL
		);
    `
    
    _, err := m.db.Exec(query)
    if err != nil {
        return fmt.Errorf("error creating Users table: %w", err)
    }
    
    fmt.Println("Users table ready")
    return nil
}