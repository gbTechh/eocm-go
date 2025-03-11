// internal/app/app.go
package app

import (
	"ecom/config"
	"ecom/internal/di"
	"ecom/internal/migrations"
	"ecom/internal/server"
	"fmt"
	"log"
)

type App struct {
    config    *config.Config
    container *di.Container
    server    *server.Server
}

func NewApp() *App {
    return &App{}
}

func (a *App) Initialize() error {
    // Cargar configuración
    cfg, err := config.LoadConfig()
    if err != nil {
        return err
    }
    a.config = cfg

    // Inicializar contenedor de dependencias
    container, err := di.NewContainer(cfg)
    if err != nil {
        return err
    }
    a.container = container

    // Inicializar servidor
    server := server.NewServer(cfg, container)
    a.server = server

		// Ejecutar migraciones después de inicializar la base de datos
    migrator := migrations.NewMigrator(a.container.DB())
    if err := migrator.Up(); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}

func (a *App) Start() error {
	defer func() {
		if err := a.container.Cleanup(); err != nil {
			log.Printf("Error during cleanup: %v", err)
		}
	}()
	return a.server.Start()
}