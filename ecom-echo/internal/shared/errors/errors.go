// internal/shared/errors/errors.go
package errors

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// Errores comunes
var (
    ErrNotFound = errors.New("not found")
)

type AppError struct {
    Code    int
    Message string
    Err     error
}


func (e *AppError) Error() string {
    return e.Message
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Code:    404,
        Message: message,
    }
}

func NewBadRequestError(message string) *AppError {
    return &AppError{
        Code:    400,
        Message: message,
    }
}

// Nueva función para errores de validación
func NewValidationError(message string) *AppError {
    return &AppError{
        Code:    400,
        Message: message,
    }
}
func NewMessageError(message string) *AppError {
    return &AppError{
        Code:    403,
        Message: message,
    }
}

// Nueva función para errores internos
func NewInternalError(message string, err error) *AppError {
    return &AppError{
        Code:    500,
        Message: message,
        Err:     err,
    }
}
func NewDuplicateEntryError(message string, err error) *AppError {
    return &AppError{
        Code:    409,
        Message: "Error campo duplicado. " + message,
        Err:     err,
    }
}
func NewMysqlError(err error) *AppError {
    
    if mysqlErr, ok := err.(*mysql.MySQLError); ok {
        // Código de error para clave duplicada en MySQL						
        if mysqlErr.Number == 1062 {
            return NewDuplicateEntryError(mysqlErr.Message, err)
        }else {
            return NewInternalError(mysqlErr.Message, err)
        }
    }		

    return NewInternalError("Error desconocido", err)
    
}


// IsNotFound verifica si un error es del tipo NotFound
func IsNotFound(err error) bool {
    if appErr, ok := err.(*AppError); ok {
        return appErr.Code == 404
    }
    return false
}