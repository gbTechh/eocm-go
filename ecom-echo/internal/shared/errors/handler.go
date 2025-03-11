// internal/shared/errors/handler.go
package errors

import (
	"ecom/internal/shared/response"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func getErrorDetail(err error) interface{} {
	if err == nil {
		return nil
	}
	if os.Getenv("APP_ENV") == "development" {
			return err.Error()
	}
	return nil
}

func HandleError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *AppError:
			return c.JSON(e.Code, response.Error(e.Message, getErrorDetail(e.Err)))
	default:
			if IsNotFound(err) {
					return c.JSON(http.StatusNotFound, 
							response.Error("Recurso no encontrado", getErrorDetail(err)))
			}
			return c.JSON(http.StatusInternalServerError, 
					response.Error("Error interno del servidor", getErrorDetail(err)))
	}
}