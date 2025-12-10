package validation

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
}

// Validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(input interface{}) error {
	if err := cv.validator.Struct(input); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (cv *CustomValidator) Validate2(input interface{}) error {
	if err := cv.validator.Struct(input); err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			fmt.Println(e)
		}

		return err
	}
	return nil
}
