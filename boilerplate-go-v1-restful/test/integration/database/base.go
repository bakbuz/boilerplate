package database_test

import (
	"testing"

	"codegen/internal/bootstrap"
	"codegen/internal/database"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

var (
	emptyGuid uuid.UUID = uuid.Nil
)

func newGuid() uuid.UUID {
	//id, _ := uuid.NewRandom()
	id, _ := uuid.NewV7()
	return id
}

func newUlid() ulid.ULID {
	id := ulid.Make()
	return id
}

func pointer[T any](v T) *T {
	return &v
}

func RequireConfig(t *testing.T) *bootstrap.Config {
	c, err := bootstrap.NewConfig("../../../config.dev.json")
	require.NoError(t, err)
	return c
}

func TestDB() (*database.DB, error) {
	db, err := database.New("./test.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
