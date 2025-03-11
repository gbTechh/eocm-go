package migrations

import (
	"database/sql"
	"ecom/internal/migrations/entities"
)

type Migration interface {
    Migrate() error
}

type Migrator struct {
    migrations []Migration
}

func NewMigrator(db *sql.DB) *Migrator {
    return &Migrator{
        migrations: []Migration{
            entities.NewUsersMigration(db),
            entities.NewRolesMigration(db),
            entities.NewModulesMigration(db),
            entities.NewAuthUserMigration(db),
            entities.NewMediaMigration(db),
            entities.NewCategoriesMigration(db),
            entities.NewTagsMigration(db),
            entities.NewTaxesMigration(db),
            entities.NewProductsMigration(db),
            entities.NewPricesMigration(db),
            entities.NewZonesMigration(db),
            entities.NewCurrenciesMigration(db),
            entities.NewCountriesMigration(db),
            entities.NewCustomersMigration(db),
            entities.NewShippingsMigration(db),
            entities.NewCartsMigration(db),
            entities.NewOrdersMigration(db),
            entities.NewTriggersMigration(db),
            //entities.NewIndexesMigration(db),
            // Agregar aquí las demás migraciones
        },
    }
}

func (m *Migrator) Up() error {
    for _, migration := range m.migrations {
        if err := migration.Migrate(); err != nil {
            return err
        }
    }
    return nil
}