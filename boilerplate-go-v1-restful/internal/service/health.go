package service

import (
	"context"

	"codegen/internal/database"
)

// HealthService ...
type HealthService interface {
	HealthCheck(ctx context.Context) error
}

type healthService struct {
	db *database.DB
}

// NewHealthService ...
func NewHealthService(db *database.DB) HealthService {
	return &healthService{db: db}
}

// HealthCheck ...
func (srv *healthService) HealthCheck(ctx context.Context) error {
	return srv.db.PingContext(ctx)
}
