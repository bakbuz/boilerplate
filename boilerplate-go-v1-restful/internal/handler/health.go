package handler

import (
	"net/http"

	"codegen/internal/service"
	"codegen/pkg"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type HealthHandler interface {
	HealthCheck(c echo.Context) error
}

type healthHandler struct {
	srv service.HealthService
}

func NewHealthHandler(srv service.HealthService) HealthHandler {
	return &healthHandler{srv: srv}
}

// HealthCheck godoc
//
//	@Summary		health check
//	@Description	health check
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	pkg.HealthResponse
//	@Failure		500	{object}	pkg.ErrorResponse
//	@Router			/health [get]
func (h *healthHandler) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()
	logger := zerolog.Ctx(ctx)

	if err := h.srv.HealthCheck(ctx); err != nil {
		logger.Error().Err(err).Msg("")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, pkg.HealthResponse{
		Status: "OK",
	})
}
