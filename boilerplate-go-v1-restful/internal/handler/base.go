package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

func shouldBindJSON(c echo.Context, req interface{}) error {
	if err := c.Bind(&req); err != nil {
		return errors.WithMessage(err, "request body is unable to bind")
	}

	if err := c.Validate(req); err != nil {
		return errors.WithMessage(err, "request body invalid")
	}
	return nil
}

// getIdFromRoute
func getIdAsInt(c echo.Context) (int, error) {
	idstr := c.Param("id")
	if idstr == "" {
		return 0, errors.New("id value is required")
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return 0, errors.WithMessage(err, "id parse error")
	}
	return id, nil
}

func getIdAsInt16(c echo.Context) (int16, error) {
	idstr := c.Param("id")
	if idstr == "" {
		return 0, errors.New("id value is required")
	}
	id, err := strconv.ParseInt(idstr, 10, 16)
	if err != nil {
		return 0, errors.WithMessage(err, "id parse error")
	}
	return int16(id), nil
}

func getIdAsUUID(c echo.Context) (uuid.UUID, error) {
	idstr := c.Param("id")
	if idstr == "" {
		return uuid.Nil, errors.New("id value is required")
	}
	id, err := uuid.Parse(idstr)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "id parse error")
	}
	return id, nil
}

func getParamAsInt16(c echo.Context, paramName string, defaultValue int16) (int16, error) {
	prmstr := c.Param(paramName)
	if prmstr == "" {
		return defaultValue, nil
	}
	prmval, err := strconv.ParseInt(prmstr, 10, 16)
	if err != nil {
		return defaultValue, errors.WithMessage(err, "failed to convert argument to integer")
	}
	value := int16(prmval)
	if value == 0 {
		value = defaultValue
	}
	return value, nil
}

// getCurrentUserId
func getCurrentUserId(c echo.Context) int {
	return 0
}
