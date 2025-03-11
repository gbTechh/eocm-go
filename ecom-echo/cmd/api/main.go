// @title Ecommerce API
// @version 1.0
// @description API para el sistema de ecommerce
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"ecom/internal/app"
	"log"
)

func main() {
    application := app.NewApp()
    // Sirve la documentación OpenAPI


    if err := application.Initialize(); err != nil {
        log.Fatal("Failed to initialize app:", err)
    }
    
    if err := application.Start(); err != nil {
        log.Fatal("Failed to start app:", err)
    }
}