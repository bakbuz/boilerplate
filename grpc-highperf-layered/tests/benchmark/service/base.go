package service_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func ptr[T any](v T) *T {
	return &v
}

type repositoryMock struct {
	mock.Mock
}

// Ping implements database.Repository
func (r *repositoryMock) Ping(ctx context.Context) error {
	args := r.Called(ctx)
	return args.Error(0)
}
