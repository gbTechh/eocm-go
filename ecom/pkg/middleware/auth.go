package middleware

import (
	"ecommerce/pkg/apperrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.NewError(
                apperrors.ErrUnauthorized,
                "Authorization header is required",
            ))
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            c.AbortWithStatusJSON(http.StatusUnauthorized, apperrors.NewError(
                apperrors.ErrUnauthorized,
                "Invalid authorization header format",
            ))
            return
        }

        // Validate token and set user in context
        // ... implementación del JWT

        c.Next()
    }
}
