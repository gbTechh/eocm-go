// pkg/response/response.go
package response

import (
	"ecommerce/pkg/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
    c.JSON(statusCode, Response{
        Success: true,
        Data:    data,
    })
}

func Error(c *gin.Context, err error) {
    var statusCode int
    var errorResponse interface{}

    switch e := err.(type) {
    case *apperrors.Error:
        statusCode = getStatusCodeForErrorCode(e.Code)
        errorResponse = e
    default:
        statusCode = http.StatusInternalServerError
        errorResponse = err.Error()
    }

    c.JSON(statusCode, Response{
        Success: false,
        Error:   errorResponse,
    })
}

func getStatusCodeForErrorCode(code apperrors.ErrorCode) int {
    switch code {
    case apperrors.ErrNotFound:
        return http.StatusNotFound
    case apperrors.ErrInvalidInput:
        return http.StatusBadRequest
    case apperrors.ErrUnauthorized:
        return http.StatusUnauthorized
    case apperrors.ErrAlreadyExists:
        return http.StatusConflict
    default:
        return http.StatusInternalServerError
    }
}