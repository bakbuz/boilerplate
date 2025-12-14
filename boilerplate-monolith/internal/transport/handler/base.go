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
		return uuid.Nil, errors.New("user context is missing")
	}

	idStr, ok := val.(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user id type")
	}

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to parse user id")
	}

	return uid, nil
}
