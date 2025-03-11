package validator

import (
	"ecommerce/pkg/apperrors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

func ValidateStruct(s interface{}) error {
    if err := validate.Struct(s); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            return apperrors.NewError(
                apperrors.ErrInvalidInput,
                "Validation failed",
                validationErrors,
            )
        }
        return err
    }
    return nil
}