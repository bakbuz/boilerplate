package handler

import (
	"codegen/internal/interceptor"
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func getUserID(ctx context.Context) (uuid.UUID, error) {
	val := ctx.Value(interceptor.UserIDKey)
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
