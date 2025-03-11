// cmd/api/main.go
package main

import (
	"context"
	"ecommerce/internal/product"
	"ecommerce/pkg/database"
	"ecommerce/pkg/logger"
	"ecommerce/pkg/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found")
    }
}

func main() {
    // Inicializar logger
    log := logger.GetLogger()
    log.Info("Starting application...")

    // Configuración de la base de datos
    dbConfig := database.Config{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
    }

    // Conectar a la base de datos
    db, err := database.NewConnection(dbConfig)
    if err != nil {
        log.Error("Failed to connect to database:", err)
        os.Exit(1)
    }
    defer db.Close()

    // Verificar conexión a la base de datos
    if err := db.Ping(); err != nil {
        log.Error("Failed to ping database:", err)
        os.Exit(1)
    }
    log.Info("Successfully connected to database")

    // Inicializar repositorios
    productRepo := product.NewRepository(db)
    /*userRepo := user.NewRepository(db)
    orderRepo := order.NewRepository(db)
    cartRepo := cart.NewRepository(db)
    authRepo := auth.NewRepository(db)*/

    // Inicializar servicios
    productService := product.NewService(productRepo)
   /* userService := user.NewService(userRepo)
    orderService := order.NewService(orderRepo, productRepo) // Note: orderService necesita productRepo para verificar stock
    cartService := cart.NewService(cartRepo, productRepo)
    authService := auth.NewService(authRepo)*/

    // Inicializar handlers
    productHandler := product.NewHandler(productService)
    /*userHandler := user.NewHandler(userService)
    orderHandler := order.NewHandler(orderService)
    cartHandler := cart.NewHandler(cartService)
    authHandler := auth.NewHandler(authService)*/

    // Configurar Gin
    gin.SetMode(os.Getenv("GIN_MODE"))
    router := gin.New()

    // Middleware global
    router.Use(gin.Recovery())
    router.Use(middleware.Logger())
    router.Use(middleware.CORS())

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "OK"})
    })

    // API versioning
    v1 := router.Group("/api/v1")

    // Rutas públicas
    {
        /*auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/forgot-password", authHandler.ForgotPassword)
            auth.POST("/reset-password", authHandler.ResetPassword)
        }*/

        products := v1.Group("/products")
        {
            products.GET("", productHandler.List)
            products.GET("/:id", productHandler.Get)
            products.GET("/categories", productHandler.ListCategories)
        }
    }

    // Rutas protegidas
    protected := v1.Group("")
    protected.Use(middleware.AuthMiddleware())
    {
        // Rutas de usuario
        /*users := protected.Group("/users")
        {
            users.GET("/me", userHandler.GetProfile)
            users.PUT("/me", userHandler.UpdateProfile)
            users.PUT("/me/password", userHandler.UpdatePassword)
        }
				*/
        // Rutas de productos (admin)
        products := protected.Group("/products")
        {
            products.POST("", middleware.AdminOnly(), productHandler.Create)
            products.PUT("/:id", middleware.AdminOnly(), productHandler.Update)
            products.DELETE("/:id", middleware.AdminOnly(), productHandler.Delete)
        }

        // Rutas de carrito
        /*cart := protected.Group("/cart")
        {
            cart.GET("", cartHandler.Get)
            cart.POST("/items", cartHandler.AddItem)
            cart.PUT("/items/:id", cartHandler.UpdateItem)
            cart.DELETE("/items/:id", cartHandler.RemoveItem)
        }*/

        // Rutas de órdenes
        /*orders := protected.Group("/orders")
        {
            orders.POST("", orderHandler.Create)
            orders.GET("", orderHandler.List)
            orders.GET("/:id", orderHandler.Get)
            orders.PUT("/:id/cancel", orderHandler.Cancel)
        }*/
    }

    // Configurar el servidor HTTP
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
        Handler: router,
    }

    // Iniciar el servidor en una goroutine
    go func() {
        log.Infof("Server starting on port %s", os.Getenv("PORT"))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Error("Failed to start server:", err)
            os.Exit(1)
        }
    }()

    // Configurar graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Error("Server forced to shutdown:", err)
    }

    log.Info("Server exiting")
}