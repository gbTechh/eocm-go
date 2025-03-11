package server

import (
	"ecom/config"
	"ecom/internal/di"
	customValidator "ecom/internal/shared/validator"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "ecom/docs"
)

type Server struct {
    echo      *echo.Echo
    config    *config.Config
    container *di.Container
}

func NewServer(cfg *config.Config, container *di.Container) *Server {
    e := echo.New()

    e.File("/docs/openapi.yaml", "docs/openapi.yaml")
    e.Static("/docs", "docs/swagger-ui")
    e.GET("/swagger/*", echoSwagger.WrapHandler)
    // Configuración básica de Echo
    setupMiddleware(e)
    setupValidator(e)
    
    return &Server{
        echo:      e,
        config:    cfg,
        container: container,
    }
}

func (s *Server) Start() error {
    // Registrar rutas
    s.registerRoutes()
    
    // Iniciar servidor
    return s.echo.Start(":" + s.config.Server.Port)
}

func setupMiddleware(e *echo.Echo) {
    e.Use(middleware.Logger())     // Log de todas las requests
    e.Use(middleware.Recover())    // Recuperación de panics
    e.Use(middleware.CORS()) 

		// Middleware de seguridad
    e.Use(middleware.Secure())
    e.Use(middleware.RequestID())
}

func setupValidator(e *echo.Echo) {
    e.Validator = customValidator.NewValidator()
}