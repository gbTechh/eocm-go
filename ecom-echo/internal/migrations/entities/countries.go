package entities

import (
	"database/sql"
	"fmt"
)

type CountriesMigration struct {
   db *sql.DB
}

func NewCountriesMigration(db *sql.DB) *CountriesMigration {
   return &CountriesMigration{db: db}
}

func (m *CountriesMigration) Migrate() error {
   // Crear tabla countries
   createCountriesQuery := `
       CREATE TABLE IF NOT EXISTS countries (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           flag VARCHAR(100) NOT NULL,
           translations JSON,
           iso2 VARCHAR(2) NOT NULL UNIQUE,
           phone_code VARCHAR(10) NOT NULL,
           id_currency INT,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           FOREIGN KEY (id_currency) REFERENCES currencies(id),
           INDEX idx_iso2 (iso2),
           INDEX idx_currency (id_currency)
       )
   `
   
   if _, err := m.db.Exec(createCountriesQuery); err != nil {
       return fmt.Errorf("error creating Countries table: %w", err)
   }

   // Crear tabla cities
   createCitiesQuery := `
       CREATE TABLE IF NOT EXISTS cities (
           id INT AUTO_INCREMENT PRIMARY KEY,
           name VARCHAR(100) NOT NULL,
           id_country INT NOT NULL,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
           FOREIGN KEY (id_country) REFERENCES countries(id),
           INDEX idx_country (id_country)
       )
   `
   
   if _, err := m.db.Exec(createCitiesQuery); err != nil {
       return fmt.Errorf("error creating Cities table: %w", err)
   }

   // Crear tabla N:M cities_zones
   createCitiesZonesQuery := `
       CREATE TABLE IF NOT EXISTS cities_zones (
           id_city INT,
           id_zone INT,
           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
           PRIMARY KEY (id_city, id_zone),
           FOREIGN KEY (id_city) REFERENCES cities(id),
           FOREIGN KEY (id_zone) REFERENCES zones(id),
           INDEX idx_zone (id_zone)
       )
   `
   
   if _, err := m.db.Exec(createCitiesZonesQuery); err != nil {
       return fmt.Errorf("error creating Cities-Zones relationship table: %w", err)
   }
   
   fmt.Println("Countries, Cities and relationships ready")
   return nil
}