package validation

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var v = validator.New()

// Validate is doing struct level validation using validator.v10
func Validate(req interface{}) error {
	// validate req message
	if err := v.Struct(req); err != nil {
		if v, ok := err.(validator.ValidationErrors); ok {
			err = v
		}
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("Request validation failed: %s", err))
	}
	return nil
}
