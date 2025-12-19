package handler

import (
	"codegen/internal/transport/interceptor"
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func getCurrentUserId(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(interceptor.UserIdKey)
	if val == nil {
		return uuid.Nil, errors.New("user_id context is missing")
	}

	userId, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("user_id context has invalid type")
	}
	return userId, nil
}
