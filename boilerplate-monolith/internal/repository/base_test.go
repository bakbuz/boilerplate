package repository_test

import (
	"codegen/internal/infrastructure/database"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// setupTestDB sets up a real database connection for integration testing.
// It requires DATABASE_URL environment variable to be set.
func setupTestDB(t *testing.T) *database.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("Skipping integration test: DATABASE_URL not set")
	}

	db, err := database.NewPool(context.Background(), dsn)
	require.NoError(t, err)

	return db
}
