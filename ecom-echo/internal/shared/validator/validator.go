// internal/shared/validator/validator.go
package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
    validator *validator.Validate
}

func NewValidator() *CustomValidator {
    v := validator.New()
    
    // Agregar validador personalizado para slug
    v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
        value := fl.Field().String()
        match, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", value)
        return match
    })
    v.RegisterValidation("name", func(fl validator.FieldLevel) bool {
        value := fl.Field().String()
        match, _ := regexp.MatchString(`^[\p{L}\s]+$`, value)
        return match
    })
    
    return &CustomValidator{
        validator: v,
    }
}

func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}