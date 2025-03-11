package middleware

import (
	"ecommerce/pkg/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obtener el usuario del contexto (establecido por AuthMiddleware)
        user, exists := c.Get("user")
        if !exists {
            c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.NewError(
                apperrors.ErrUnauthorized,
                "Usuario no autenticado",
            ))
            return
        }

        // Verificar si el usuario es admin
        // Asumiendo que el usuario tiene un campo Role
        if user.(map[string]interface{})["role"] != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, apperrors.NewError(
                apperrors.ErrUnauthorized,
                "Se requieren privilegios de administrador",
            ))
            return
        }

        c.Next()
    }
}