package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

func LoadConfig() (*Config, error) {
	// Cargar .env
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{}
	
	// Server
	config.Server.Port = os.Getenv("PORT")
	if config.Server.Port == "" {
		config.Server.Port = "8080" // valor por defecto
	}

	// Database
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Name = os.Getenv("DB_NAME")

	// Validaciones
	if config.Database.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Port == "" {
		return nil, fmt.Errorf("DB_PORT is required")
	}
	if config.Database.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.Database.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return config, nil
}